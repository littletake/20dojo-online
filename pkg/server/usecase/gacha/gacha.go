package gacha

// import (
// 	"fmt"
// 	"log"

// 	"20dojo-online/pkg/constant"
// 	"20dojo-online/pkg/db"
// 	"20dojo-online/pkg/server/domain/model"
// 	"20dojo-online/pkg/server/domain/repository"
// 	"20dojo-online/pkg/server/usecase"
// )

// // GachaUseCase UserにおけるUseCaseのインターフェース
// type GachaUseCase interface {
// 	Gacha(gachaTimes int32, userID string) ([]*model.GachaResult, *model.MyErr)
// 	CreateItemRatioSlice() ([]*model.GachaProb, *model.MyErr)
// }

// type gachaUseCase struct {
// 	userRepository      repository.UserRepository
// 	cItemRepository     repository.CItemRepository
// 	ucItemRepository    repository.UCItemRepository
// 	gachaProbRepository repository.GachaProbRepository
// }

// // NewGachaUseCase Userデータに関するUseCaseを生成
// func NewGachaUseCase(ur repository.UserRepository, cr repository.CItemRepository,
// 	ucr repository.UCItemRepository, gpr repository.GachaProbRepository) GachaUseCase {
// 	return &gachaUseCase{
// 		userRepository:      ur,
// 		cItemRepository:     cr,
// 		ucItemRepository:    ucr,
// 		gachaProbRepository: gpr,
// 	}
// }

// // GetUsersByHighScore Userデータを条件抽出
// func (gu gachaUseCase) Gacha(gachaTimes int32, userID string) ([]*model.GachaResult, *model.MyErr) {
// 	// userIDと照合するユーザを取得
// 	// idと照合するユーザを取得
// 	user, err := gu.userRepository.SelectUserByUserID(userID)
// 	if err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	if user == nil {
// 		myErr := usecase.CreateMyErr(
// 			fmt.Errorf("user not found"),
// 			500,
// 		)
// 		return nil, myErr
// 	}
// 	// 必要枚数分のコインがあるかどうかを判定
// 	necessaryCoins := constant.GachaCoinConsumption * gachaTimes
// 	if user.Coin-necessaryCoins < 0 {
// 		myErr := usecase.CreateMyErr(
// 			fmt.Errorf("user doesn't have enough coins. current: %d, necessary: %d", user.Coin, necessaryCoins),
// 			400,
// 		)
// 		return nil, myErr
// 	}
// 	// TODO: rollbackの時に対応できない！
// 	user.Coin = user.Coin - necessaryCoins

// 	// table: collection_itemの全件取得
// 	cItemSlice, err := gu.cItemRepository.SelectAllCollectionItem()
// 	if err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	// table: user_collection_itemに対してuserIDのものを取得
// 	ucItemSlice, err := gu.ucItemRepository.SelectUCItemSliceByUserID(userID)
// 	if err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	// hasGotItemMap 既出アイテム一覧map
// 	// [注意] ガチャ実行時も追加するので可変長指定
// 	hasGotItemMap := make(map[string]bool)
// 	for _, ucItem := range ucItemSlice {
// 		itemID := ucItem.CollectionItemID
// 		hasGotItemMap[itemID] = true
// 	}

// 	// 1. 乱数によるガチャの実行
// 	// gettingItemSlice 当てたアイテムのIDのslice
// 	gettingItemSlice, myErr := GetItemSlice(gu, gachaTimes, cItemSlice)
// 	if myErr != nil {
// 		return nil, myErr
// 	}

// 	// 2. アイテムの照合
// 	// TODO: アイテムの保存
// 	gachaResultSlice := make([]*model.GachaResult, gachaTimes)
// 	newItemSlice := ChangeToGachaResult(gettingItemSlice, cItemSlice, hasGotItemMap, gachaResultSlice, userID)
// 	// for i, gettingItem := range gettingItemSlice {
// 	// 	for _, item := range cItemSlice {
// 	// 		if gettingItem == item.ItemID {
// 	// 			// 既出itemの確認
// 	// 			if hasGotItemMap[item.ItemID] {
// 	// 				result := model.GachaResult{
// 	// 					CollectionID: item.ItemID,
// 	// 					ItemName:     item.ItemName,
// 	// 					Rarity:       item.Rarity,
// 	// 					IsNew:        false,
// 	// 				}
// 	// 				gachaResultSlice[i] = &result
// 	// 			} else {
// 	// 				result := model.GachaResult{
// 	// 					CollectionID: item.ItemID,
// 	// 					ItemName:     item.ItemName,
// 	// 					Rarity:       item.Rarity,
// 	// 					IsNew:        true,
// 	// 				}
// 	// 				gachaResultSlice[i] = &result
// 	// 				// 既出アイテム一覧に追加
// 	// 				hasGotItemMap[item.ItemID] = true
// 	// 				// 登録
// 	// 				newItem := model.UserCollectionItem{
// 	// 					UserID:           userID,
// 	// 					CollectionItemID: item.ItemID,
// 	// 				}
// 	// 				newItemSlice = append(newItemSlice, &newItem)
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	// TODO: トランザクションのテスト作成
// 	// 3. トランザクション開始（複数DB操作）
// 	tx, err := db.Conn.Begin()
// 	if err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	// TODO: 書き方再検討
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Println("!! PANIC !!")
// 			log.Println(err)
// 			if rollbackErr := tx.Rollback(); rollbackErr != nil {
// 				log.Println("failed to Rollback")
// 				log.Println(rollbackErr)
// 				// myErr := CreateMyErr(rollbackErr, 500)
// 				// return nil, myErr
// 			}
// 		}
// 		// TODO: ここにpanic()を書くとAPIが二回実行される
// 		// panic(err)
// 	}()
// 	// 3-1. バルクインサート
// 	if len(newItemSlice) != 0 {
// 		if err := gu.ucItemRepository.BulkInsertUCItemSlice(newItemSlice, tx); err != nil {
// 			myErr := usecase.CreateMyErr(err, 500)
// 			return nil, myErr
// 		}
// 	}
// 	// 3-2. ユーザの保持コイン更新
// 	if err := gu.userRepository.UpdateUserByUserInTx(user, tx); err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	if err := tx.Commit(); err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	return gachaResultSlice, nil
// }

// func (gu gachaUseCase) CreateItemRatioSlice() ([]*model.GachaProb, *model.MyErr) {
// 	gachaProbSlice, err := gu.gachaProbRepository.SelectAllGachaProb()
// 	if err != nil {
// 		myErr := usecase.CreateMyErr(err, 500)
// 		return nil, myErr
// 	}
// 	return gachaProbSlice, nil
// }
