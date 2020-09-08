package gacha

import (
	"database/sql"
	"fmt"
	"math/rand"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/server/domain/model"
	cr "20dojo-online/pkg/server/domain/repository/collection_item"
	gr "20dojo-online/pkg/server/domain/repository/gacha_probability"
	txr "20dojo-online/pkg/server/domain/repository/transaction"
	ur "20dojo-online/pkg/server/domain/repository/user"
	ucr "20dojo-online/pkg/server/domain/repository/user_collection_item"
	"20dojo-online/pkg/server/interface/myerror"
)

// GachaUseCase UserにおけるUseCaseのインターフェース
type GachaUseCase interface {
	Gacha(gachaTimes int32, userID string) ([]*model.GachaResult, *myerror.MyErr)
	CreateItemRatioSlice() *myerror.MyErr
	CreateCItemSlice() *myerror.MyErr
	GetItems(gachaTimes int32) []string
	CreateGachaResults(gettingItemSlice []string, hasGotItemMap map[string]bool, userID string) ([]*model.GachaResult, []*model.UserCollectionItem)
}

type gachaUseCase struct {
	userRepository      ur.UserRepository
	cItemRepository     cr.CItemRepository
	ucItemRepository    ucr.UCItemRepository
	gachaProbRepository gr.GachaProbRepository
	seed                int64
	txRepository        txr.TxRepository
}

// NewGachaUseCase Userデータに関するUseCaseを生成
func NewGachaUseCase(ur ur.UserRepository, cr cr.CItemRepository,
	ucr ucr.UCItemRepository, gpr gr.GachaProbRepository, seed int64, txr txr.TxRepository) GachaUseCase {
	return &gachaUseCase{
		userRepository:      ur,
		cItemRepository:     cr,
		ucItemRepository:    ucr,
		gachaProbRepository: gpr,
		seed:                seed,
		txRepository:        txr,
	}
}

// // init 乱数のseed定義
// func init() {
// }

// cItemSlice collectionItemのスライス
var cItemSlice []*model.CollectionItem

// hasGotcItemSlice table:collection_itemの取得状況
var hasGotcItemSlice bool

// itemRatioSlice ratioを考慮したアイテム対応表
var itemRatioSlice []int32

// hasGotGachaProb table:gacha_probabilityの取得状況
var hasGotGachaProb bool

// GetUsersByHighScore Userデータを条件抽出
func (gu gachaUseCase) Gacha(gachaTimes int32, userID string) ([]*model.GachaResult, *myerror.MyErr) {
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
	ucItemSlice, err := gu.ucItemRepository.SelectUCItemSliceByUserID(userID)
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

	// TODO: できていない（結局txが必要になってしまう）
	// 3. トランザクション開始（複数DB操作）
	if err = gu.txRepository.Transaction(gu.BulkInsertAndUpdate(newItemSlice, user, tx)); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}

	return gachaResultSlice, nil
}

func (gu gachaUseCase) CreateCItemSlice() *myerror.MyErr {
	// 全アイテムの情報をまとめる -> cItemSlice
	// hasGotItemSlice 一度だけ実行するように制御
	var err error
	if !hasGotcItemSlice {
		cItemSlice, err = gu.cItemRepository.SelectAllCollectionItem()
		if err != nil {
			myErr := myerror.NewMyErr(err, 500)
			return myErr
		}
		hasGotcItemSlice = true
		return nil
	}
	return nil
}

// CreateItemRatioSlice ratioを考慮したアイテム対応表の作成
func (gu gachaUseCase) CreateItemRatioSlice() *myerror.MyErr {
	// gacha_probabilityの情報を取得
	if !hasGotGachaProb {
		gachaProbSlice, err := gu.gachaProbRepository.SelectAllGachaProb()
		if err != nil {
			myErr := myerror.NewMyErr(err, 500)
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
func (gu gachaUseCase) GetItems(gachaTimes int32) []string {
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
func (gu gachaUseCase) CreateGachaResults(gettingItemSlice []string, hasGotItemMap map[string]bool, userID string) ([]*model.GachaResult, []*model.UserCollectionItem) {
	gachaResultSlice := make([]*model.GachaResult, len(gettingItemSlice))
	newItemSlice := []*model.UserCollectionItem{}
	for i, gettingItem := range gettingItemSlice {
		for _, item := range cItemSlice {
			if gettingItem == item.ItemID {
				// レスポンス用に整形
				result := model.GachaResult{
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
					newItem := model.UserCollectionItem{
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

// TODO: 未完成
func (gu gachaUseCase) BulkInsertAndUpdate(newItemSlice []*model.UserCollectionItem, user *model.UserL, tx *sql.Tx) error {
	// 3-1. バルクインサート
	if len(newItemSlice) != 0 {
		if err := gu.ucItemRepository.BulkInsertUCItemSlice(newItemSlice, tx); err != nil {
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
