package initializer

// import (
// 	"20dojo-online/pkg/server/domain/model"
// 	"20dojo-online/pkg/server/usecase"
// )

// // ItemRatioSlice ratioを考慮したアイテム対応表
// var ItemRatioSlice []int32

// // CreateItemRatioSliceOnce アイテム対応表の作成
// func CreateItemRatioSliceOnce(gachaUseCase usecase.GachaUseCase) *model.MyErr {
// 	// table:gacha_probabilityの全件取得
// 	gachaProbSlice, myErr := gachaUseCase.CreateItemRatioSlice()
// 	if myErr != nil {
// 		return myErr
// 	}
// 	ItemRatioSlice = make([]int32, len(gachaProbSlice))
// 	count := int32(0)
// 	for i, item := range gachaProbSlice {
// 		count += item.Ratio
// 		ItemRatioSlice[i] = count
// 	}
// 	return nil
// }
