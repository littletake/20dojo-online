package usecase

import "20dojo-online/pkg/server/domain/model"

// グローバル変数を明記
// collectionItemのスライス
var cItemSlice []*model.CollectionItem

// Collection_itemの情報の取得状況
var hasGotcItemSlice bool

// itemRatioSlice ratioを考慮したアイテム対応表
var itemRatioSlice []int32

// gacha_probabilityの情報の取得状況
var hasGotGachaProb bool
