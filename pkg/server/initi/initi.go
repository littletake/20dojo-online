package initi

import (
	"20dojo-online/pkg/server/model"
)

// itemRatioSlice ratioを考慮したアイテム対応表
var itemRatioSlice []int32

// createItemRatioSliceOnce() アイテム対応表の作成
func createItemRatioSliceOnce() error {
	// table:gacha_probabilityの全件取得
	gachaProbSlice, err := model.SelectAllGachaProb()
	if err != nil {
		return err
	}
	itemRatioSlice = make([]int32, len(gachaProbSlice))
	count := int32(0)
	for i, item := range gachaProbSlice {
		count += item.Ratio
		itemRatioSlice[i] = count
	}
	return nil
}
