//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package gacha

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"20dojo-online/pkg/constant"
	cm "20dojo-online/pkg/server/domain/model/collectionitem"
	um "20dojo-online/pkg/server/domain/model/user"
	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
	cr "20dojo-online/pkg/server/domain/repository/collectionitem"
	gr "20dojo-online/pkg/server/domain/repository/gachaprobability"

	txr "20dojo-online/pkg/server/domain/repository/transaction"
	ur "20dojo-online/pkg/server/domain/repository/user"
	ucr "20dojo-online/pkg/server/domain/repository/usercollectionitem"
	"20dojo-online/pkg/server/interface/myerror"
)

// GachaUseCase UserにおけるUseCaseのインターフェース
type GachaUseCase interface {
	Gacha(gachaTimes int32, userID string) ([]*GachaResult, *myerror.MyErr)
	CreateItemRatioSlice() *myerror.MyErr
	CreateCItemSlice() *myerror.MyErr
	GetItems(gachaTimes int32) []string
	CreateGachaResults(gettingItemSlice []string, hasGotItemMap map[string]bool, userID string) ([]*GachaResult, []*ucm.UserCollectionItem)
	BulkInsertAndUpdate(newItemSlice []*ucm.UserCollectionItem, user *um.UserL, tx *sql.Tx) error
}

type gachaUseCase struct {
	userRepository      ur.UserRepo
	cItemRepository     cr.CollectionItemRepo
	ucItemRepository    ucr.UserCollectionItemRepo
	gachaProbRepository gr.GachaProbRepo
	seed                int64
	txRepository        txr.TxRepo
}

// NewGachaUseCase Userデータに関するUseCaseを生成
func NewGachaUseCase(ur ur.UserRepo, cr cr.CollectionItemRepo,
	ucr ucr.UserCollectionItemRepo, gpr gr.GachaProbRepo, seed int64, txr txr.TxRepo) GachaUseCase {
	return &gachaUseCase{
		userRepository:      ur,
		cItemRepository:     cr,
		ucItemRepository:    ucr,
		gachaProbRepository: gpr,
		seed:                seed,
		txRepository:        txr,
	}
}

// GachaResult レスポンス用の構造体
type GachaResult struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

// cItemSlice collectionItemのスライス
var cItemSlice []*cm.CollectionItem

// hasGotcItemSlice table:collection_itemの取得状況
var hasGotcItemSlice bool

// itemRatioSlice ratioを考慮したアイテム対応表
var itemRatioSlice []int32

// hasGotGachaProb table:gacha_probabilityの取得状況
var hasGotGachaProb bool

// GetUsersByHighScore Userデータを条件抽出
func (gu *gachaUseCase) Gacha(gachaTimes int32, userID string) ([]*GachaResult, *myerror.MyErr) {
	// userIDと照合するユーザを取得
	user, err := gu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found"),
			500,
		)
		return nil, myErr
	}
	// 必要枚数分のコインがあるかどうかを判定
	necessaryCoins := constant.GachaCoinConsumption * gachaTimes
	if user.Coin-necessaryCoins < 0 {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user doesn't have enough coins. actual: %d, expected: %d", user.Coin, necessaryCoins),
			400,
		)
		return nil, myErr
	}
	user.Coin = user.Coin - necessaryCoins

	// --- DBの情報を取得してガチャの前準備をする

	// 現ユーザが保持しているアイテムの情報をまとめる -> hasGotItemMap
	// table: user_collection_itemに対してuserIDのものを取得
	ucItemSlice, err := gu.ucItemRepository.SelectSliceByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	// hasGotItemMap 既出アイテム一覧map
	// [注意] ガチャ実行時も追加するので可変長指定
	hasGotItemMap := make(map[string]bool)
	for _, ucItem := range ucItemSlice {
		itemID := ucItem.CollectionItemID
		hasGotItemMap[itemID] = true
	}
	// 全アイテムの情報をまとめる -> cItemSlice
	// hasGotItemSlice 一度だけ実行するように制御
	if myErr := gu.CreateCItemSlice(); myErr != nil {
		return nil, myErr
	}
	// ratioを考慮したアイテム対応表の作成
	// hasGotGachaProb 一度だけ実行するように制御
	if myErr := gu.CreateItemRatioSlice(); myErr != nil {
		return nil, myErr
	}

	// --- ガチャの実行
	// 1. 乱数によるアイテムの取得 -> 実行結果: gettingItemSlice
	// gettingItemSlice 当てたアイテムのIDのslice
	gettingItemSlice := gu.GetItems(gachaTimes)

	// 2. アイテムの照合
	gachaResultSlice, newItemSlice := gu.CreateGachaResults(gettingItemSlice, hasGotItemMap, userID)

	// // 3. トランザクション開始（複数DB操作）
	if err = gu.txRepository.Transaction(func(tx *sql.Tx) error {
		err := gu.BulkInsertAndUpdate(newItemSlice, user, tx)
		return err
	}); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}

	return gachaResultSlice, nil
}

// CreateCItemSlice collectionItemSliceを作成
func (gu *gachaUseCase) CreateCItemSlice() *myerror.MyErr {
	// 全アイテムの情報をまとめる -> cItemSlice
	// hasGotItemSlice 一度だけ実行するように制御
	var err error
	if !hasGotcItemSlice {
		cItemSlice, err = gu.cItemRepository.SelectAllCollectionItem()
		if err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}
		hasGotcItemSlice = true
		return nil
	}
	return nil
}

// CreateItemRatioSlice ratioを考慮したアイテム対応表の作成
func (gu *gachaUseCase) CreateItemRatioSlice() *myerror.MyErr {
	// gacha_probabilityの情報を取得
	if !hasGotGachaProb {
		gachaProbSlice, err := gu.gachaProbRepository.SelectAllGachaProb()
		if err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}
		itemRatioSlice = make([]int32, len(gachaProbSlice))
		count := int32(0)
		for i, item := range gachaProbSlice {
			count += item.Ratio
			itemRatioSlice[i] = count
		}
		hasGotGachaProb = true
	}
	return nil
}

// GetItems 乱数によるアイテムの取得 -> 実行結果: gettingItemSlice
func (gu *gachaUseCase) GetItems(gachaTimes int32) []string {
	// gettingItemSlice 当てたアイテムのIDのslice
	rand.Seed(gu.seed)
	gettingItemSlice := make([]string, gachaTimes)
	for i := int32(0); i < gachaTimes; i++ {
		randomNum := rand.Int31n(itemRatioSlice[len(itemRatioSlice)-1])
		index := detectNumber(randomNum)
		gettingItemSlice[i] = cItemSlice[index].ItemID
	}
	return gettingItemSlice
}

// CreateGachaResults ガチャ実行結果の作成
func (gu *gachaUseCase) CreateGachaResults(gettingItemSlice []string, hasGotItemMap map[string]bool, userID string) ([]*GachaResult, []*ucm.UserCollectionItem) {
	gachaResultSlice := make([]*GachaResult, len(gettingItemSlice))
	newItemSlice := []*ucm.UserCollectionItem{}
	for i, gettingItem := range gettingItemSlice {
		for _, item := range cItemSlice {
			if gettingItem == item.ItemID {
				// レスポンス用に整形
				result := GachaResult{
					CollectionID: item.ItemID,
					ItemName:     item.ItemName,
					Rarity:       item.Rarity,
					IsNew:        !hasGotItemMap[item.ItemID],
				}
				gachaResultSlice[i] = &result
				// 既出itemの確認
				if !hasGotItemMap[item.ItemID] {
					// 既出アイテム一覧に追加
					hasGotItemMap[item.ItemID] = true
					newItem := ucm.UserCollectionItem{
						UserID:           userID,
						CollectionItemID: item.ItemID,
					}
					newItemSlice = append(newItemSlice, &newItem)
				}
			}
		}
	}
	return gachaResultSlice, newItemSlice

}

func (gu *gachaUseCase) BulkInsertAndUpdate(newItemSlice []*ucm.UserCollectionItem, user *um.UserL, tx *sql.Tx) error {
	// 3-1. バルクインサート
	if len(newItemSlice) != 0 {
		if err := gu.ucItemRepository.BulkInsert(newItemSlice, tx); err != nil {
			return err
		}
	}
	// 3-2. ユーザの保持コイン更新
	if err := gu.userRepository.UpdateUserByUserInTx(user, tx); err != nil {
		return err
	}
	return nil
}

// TODO: 当たっているかどうかを判定する関数を作成すること
// detectNumber 適している番号を見つける
func detectNumber(random int32) int32 {
	num := int32(0)
	for {
		if itemRatioSlice[num] > random {
			break
		}
		num++
	}
	return num
}
