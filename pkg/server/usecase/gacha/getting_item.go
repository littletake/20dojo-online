package gacha

// import (
// 	"math/rand"

// 	"20dojo-online/pkg/server/domain/model"
// 	"20dojo-online/pkg/server/usecase"
// )

// // itemRatioSlice ratioを考慮したアイテム対応表
// var itemRatioSlice []int32

// // アイテム対応表の存在の有無
// var hasGotItemRatioSlice bool

// // GetItemSlice 乱数によるアイテムの取得
// // TODO: ItemRatioSlice をredisで実現
// func GetItemSlice(gachaUseCase GachaUseCase, gachaTimes int32, cItemSlice []*model.CollectionItem) ([]string, *model.MyErr) {
// 	// 1. アイテム対応表の作成
// 	if myErr := createItemRatioSlice(gachaUseCase); myErr != nil {
// 		return nil, myErr
// 	}
// 	// 2. 当てたアイテムのIDのslice作成
// 	// itemSlice 当てたアイテムのIDのslice
// 	itemSlice := make([]string, gachaTimes)
// 	for i := int32(0); i < gachaTimes; i++ {
// 		randomNum := rand.Int31n(itemRatioSlice[len(itemRatioSlice)-1])
// 		index := detectNumber(randomNum)
// 		itemSlice[i] = cItemSlice[index].ItemID
// 	}
// 	return itemSlice, nil
// }

// // TODO: テストすること
// // CreateItemRatioSliceOnce アイテム対応表の作成
// func createItemRatioSlice(gachaUseCase usecase.GachaUseCase) *model.MyErr {
// 	if !hasGotItemRatioSlice {
// 		// table:gacha_probabilityの全件取得
// 		gachaProbSlice, myErr := gachaUseCase.CreateItemRatioSlice()
// 		if myErr != nil {
// 			return myErr
// 		}
// 		itemRatioSlice = make([]int32, len(gachaProbSlice))
// 		count := int32(0)
// 		for i, item := range gachaProbSlice {
// 			count += item.Ratio
// 			itemRatioSlice[i] = count
// 		}
// 	}
// 	return nil
// }

// // detectNumber 適している番号を見つける
// func detectNumber(random int32) int32 {
// 	// TODO: 当たっているかどうかを判定する関数を作成すること
// 	num := int32(0)
// 	for {
// 		if itemRatioSlice[num] > random {
// 			break
// 		}
// 		num++
// 	}
// 	return num
// }
