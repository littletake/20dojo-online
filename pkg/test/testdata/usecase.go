package testdata

import (
	cm "20dojo-online/pkg/server/domain/model/collectionitem"
	gm "20dojo-online/pkg/server/domain/model/gachaprobability"
	um "20dojo-online/pkg/server/domain/model/user"
	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
	crm "20dojo-online/pkg/server/usecase/collection"
	grm "20dojo-online/pkg/server/usecase/gacha"
)

// ExampleUser userL の例
var ExampleUser = &um.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 0,
	Coin:      0,
}

// ExampleCItemResult1 CollectionItemResult の例
var ExampleCItemResult1 = &crm.CollectionItemResult{
	CollectionID: exampleCItem1.ItemID,
	ItemName:     exampleCItem1.ItemName,
	Rarity:       exampleCItem1.Rarity,
	HasItem:      true,
}

// ExampleCItemResult2 CollectionItemResult の例
var ExampleCItemResult2 = &crm.CollectionItemResult{
	CollectionID: exampleCItem2.ItemID,
	ItemName:     exampleCItem2.ItemName,
	Rarity:       exampleCItem2.Rarity,
	HasItem:      false,
}

// ExampleCItemResult3 CollectionItemResult の例
var ExampleCItemResult3 = &crm.CollectionItemResult{
	CollectionID: exampleCItem3.ItemID,
	ItemName:     exampleCItem3.ItemName,
	Rarity:       exampleCItem3.Rarity,
	HasItem:      false,
}

// ReturnUCItemSlice user_collection_item の例
var ReturnUCItemSlice = []*ucm.UserCollectionItem{
	exampleUCItem1,
}
var exampleUCItem1 = &ucm.UserCollectionItem{
	UserID:           ExampleUser.ID,
	CollectionItemID: exampleCItem1.ItemID,
}

// ReturnCItemSlice collection_item の例
var ReturnCItemSlice = []*cm.CollectionItem{
	exampleCItem1,
	exampleCItem2,
	exampleCItem3,
}
var exampleCItem1 = &cm.CollectionItem{
	ItemID:   "1001",
	ItemName: "example1",
	Rarity:   int32(1),
}
var exampleCItem2 = &cm.CollectionItem{
	ItemID:   "1002",
	ItemName: "example2",
	Rarity:   int32(2),
}
var exampleCItem3 = &cm.CollectionItem{
	ItemID:   "1003",
	ItemName: "example3",
	Rarity:   int32(3),
}

// // ExampleGachaResultSlice gacha_probabiliy の例
// var ExampleGachaResultSlice = []*model.GachaResult{
// 	exampleGachaResult1,
// 	exampleGachaResult2,
// 	exampleGachaResult3,
// }

// ExampleGachaResult1 GachaResult の例
var ExampleGachaResult1 = &grm.GachaResult{
	CollectionID: "1001",
	ItemName:     "example1",
	Rarity:       int32(1),
	IsNew:        false,
}

// ExampleGachaResult2 GachaResult の例
var ExampleGachaResult2 = &grm.GachaResult{
	CollectionID: "1002",
	ItemName:     "example2",
	Rarity:       int32(2),
	IsNew:        false,
}

// ExampleGachaResult3 GachaResult の例
var ExampleGachaResult3 = &grm.GachaResult{
	CollectionID: "1003",
	ItemName:     "example3",
	Rarity:       int32(3),
	IsNew:        false,
}

// ReturnGachaProbSlice GachaProb のスライスの例
var ReturnGachaProbSlice = []*gm.GachaProb{
	exampleGachaProb1,
	exampleGachaProb2,
	exampleGachaProb3,
}
var exampleGachaProb1 = &gm.GachaProb{
	CollectionItemID: "1001",
	Ratio:            6,
}
var exampleGachaProb2 = &gm.GachaProb{
	CollectionItemID: "1002",
	Ratio:            6,
}
var exampleGachaProb3 = &gm.GachaProb{
	CollectionItemID: "1003",
	Ratio:            6,
}
