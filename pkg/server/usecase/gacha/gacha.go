package gacha

import (
	"fmt"
	"log"
	"math/rand"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
	"20dojo-online/pkg/server/interface/myerror"
)

// GachaUseCase UserにおけるUseCaseのインターフェース
type GachaUseCase interface {
	Gacha(gachaTimes int32, userID string) ([]*model.GachaResult, *myerror.MyErr)
	CreateItemRatioSlice() *myerror.MyErr
}

type gachaUseCase struct {
	userRepository      repository.UserRepository
	cItemRepository     repository.CItemRepository
	ucItemRepository    repository.UCItemRepository
	gachaProbRepository repository.GachaProbRepository
}

// NewGachaUseCase Userデータに関するUseCaseを生成
func NewGachaUseCase(ur repository.UserRepository, cr repository.CItemRepository,
	ucr repository.UCItemRepository, gpr repository.GachaProbRepository) GachaUseCase {
	return &gachaUseCase{
		userRepository:      ur,
		cItemRepository:     cr,
		ucItemRepository:    ucr,
		gachaProbRepository: gpr,
	}
}

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
	if !hasGotcItemSlice {
		cItemSlice, err = gu.cItemRepository.SelectAllCollectionItem()
		if err != nil {
			myErr := myerror.NewMyErr(err, 500)
			return nil, myErr
		}
		hasGotcItemSlice = true
	}
	// ratioを考慮したアイテム対応表の作成
	// hasGotGachaProb 一度だけ実行するように制御
	if !hasGotGachaProb {
		if myErr := gu.CreateItemRatioSlice(); myErr != nil {
			return nil, myErr
		}
		hasGotGachaProb = true
	}

	// --- ガチャの実行
	// 1. 乱数によるアイテムの取得 -> 実行結果: gettingItemSlice
	// gettingItemSlice 当てたアイテムのIDのslice
	gettingItemSlice := make([]string, gachaTimes)
	for i := int32(0); i < gachaTimes; i++ {
		randomNum := rand.Int31n(itemRatioSlice[len(itemRatioSlice)-1])
		index := detectNumber(randomNum)
		// TODO: 以下のitemIDを取得する部分，collectionItemSliceではなくgachaProbで代用できない??
		gettingItemSlice[i] = cItemSlice[index].ItemID
	}
	// 2. アイテムの照合
	gachaResultSlice := make([]*model.GachaResult, gachaTimes)
	var newItemSlice []*model.UserCollectionItem
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
	// TODO: トランザクションのテスト作成
	// 3. トランザクション開始（複数DB操作）
	tx, err := db.Conn.Begin()
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	// TODO: 書き方再検討
	defer func() {
		if err := recover(); err != nil {
			log.Println("!! PANIC !!")
			log.Println(err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("failed to Rollback")
				log.Println(rollbackErr)
				// myErr := CreateMyErr(rollbackErr, 500)
				// return nil, myErr
			}
		}
		// TODO: ここにpanic()を書くとAPIが二回実行される
		// panic(err)
	}()
	// 3-1. バルクインサート
	if len(newItemSlice) != 0 {
		if err := gu.ucItemRepository.BulkInsertUCItemSlice(newItemSlice, tx); err != nil {
			myErr := myerror.NewMyErr(err, 500)
			return nil, myErr
		}
	}
	// 3-2. ユーザの保持コイン更新
	if err := gu.userRepository.UpdateUserByUserInTx(user, tx); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	if err := tx.Commit(); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	return gachaResultSlice, nil
}

// CreateItemRatioSlice ratioを考慮したアイテム対応表の作成
func (gu gachaUseCase) CreateItemRatioSlice() *myerror.MyErr {
	// gacha_probabilityの情報を取得
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

// // TODO: 役割ごとにメソッドを分けたい（分けようとして失敗）
// // getItemSlice アイテム取得
// func (gu gachaUseCase) getItemSlice(gachaTimes int32) ([]string, *myerror.MyErr) {
// 	// 1. ratioを考慮したアイテム対応表の作成
// 	var gachaProbSlice []*model.GachaProb
// 	if !hasGotGachaProb {
// 		// gacha_probabilityの情報を取得
// 		gachaProbSlice, err := gu.gachaProbRepository.SelectAllGachaProb()
// 		if err != nil {
// 			myErr := myerror.MyErr{err, 500}
// 			return nil, &myErr
// 		}
// 		itemRatioSlice = make([]int32, len(gachaProbSlice))
// 		count := int32(0)
// 		for i, item := range gachaProbSlice {
// 			count += item.Ratio
// 			itemRatioSlice[i] = count
// 		}
// 		hasGotGachaProb = true
// 	}
// 	// 2. 乱数生成して対応表からitemIDを特定
// 	gettingItemSlice := make([]string, gachaTimes)
// 	for i := int32(0); i < gachaTimes; i++ {
// 		randomNum := rand.Int31n(itemRatioSlice[len(itemRatioSlice)-1])
// 		index := detectNumber(randomNum)
// 		fmt.Print(gachaProbSlice)
// 		fmt.Print(index)
// 		gettingItemSlice[i] = gachaProbSlice[index].CollectionItemID
// 	}
// 	return gettingItemSlice, nil
// }
