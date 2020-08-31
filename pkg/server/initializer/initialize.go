package initializer

import (
	"20dojo-online/pkg/server/model"
)

// ItemRatioSlice ratioを考慮したアイテム対応表
var ItemRatioSlice []int32

// CreateItemRatioSliceOnce アイテム対応表の作成
func CreateItemRatioSliceOnce() error {
	// table:gacha_probabilityの全件取得
	gachaProbSlice, err := model.SelectAllGachaProb()
	if err != nil {
		return err
	}
	ItemRatioSlice = make([]int32, len(gachaProbSlice))
	count := int32(0)
	for i, item := range gachaProbSlice {
		count += item.Ratio
		ItemRatioSlice[i] = count
	}
	return nil
}
