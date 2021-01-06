package gacha

import (
	cm "20dojo-online/pkg/server/domain/model/collectionitem"
	gm "20dojo-online/pkg/server/domain/model/gachaprobability"
	um "20dojo-online/pkg/server/domain/model/user"
	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
)

// ExampleUser UserLモデルの例
var ExampleUser = &um.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 100,
	Coin:      100,
}

var ExampleCItemSlice = []*cm.CollectionItem{
	ExampleCItem1,
	ExampleCItem2,
	ExampleCItem3,
}
var ExampleCItem1 = &cm.CollectionItem{
	ItemID:   "1001",
	ItemName: "example1",
	Rarity:   int32(1),
}
var ExampleCItem2 = &cm.CollectionItem{
	ItemID:   "1002",
	ItemName: "example2",
	Rarity:   int32(2),
}
var ExampleCItem3 = &cm.CollectionItem{
	ItemID:   "1003",
	ItemName: "example3",
	Rarity:   int32(3),
}

// var ExampleCItemResult1 = &crm.CollectionItemResult{
// 	CollectionID: ExampleCItem1.ItemID,
// 	ItemName:     ExampleCItem1.ItemName,
// 	Rarity:       ExampleCItem1.Rarity,
// 	HasItem:      true,
// }
// var ExampleCItemResult2 = &crm.CollectionItemResult{
// 	CollectionID: ExampleCItem2.ItemID,
// 	ItemName:     ExampleCItem2.ItemName,
// 	Rarity:       ExampleCItem2.Rarity,
// 	HasItem:      false,
// }
// var ExampleCItemResult3 = &crm.CollectionItemResult{
// 	CollectionID: ExampleCItem3.ItemID,
// 	ItemName:     ExampleCItem3.ItemName,
// 	Rarity:       ExampleCItem3.Rarity,
// 	HasItem:      false,
// }

var ExampleUCItemSlice = []*ucm.UserCollectionItem{
	exampleUCItem1,
}
var exampleUCItem1 = &ucm.UserCollectionItem{
	UserID:           ExampleUser.ID,
	CollectionItemID: ExampleCItem1.ItemID,
}

var NewItem1 = &ucm.UserCollectionItem{
	UserID:           ExampleUser.ID,
	CollectionItemID: ExampleCItem2.ItemID,
}
var NewItem2 = &ucm.UserCollectionItem{
	UserID:           ExampleUser.ID,
	CollectionItemID: ExampleCItem3.ItemID,
}

var ExampleGachaResult1 = &GachaResult{
	CollectionID: "1001",
	ItemName:     "example1",
	Rarity:       int32(1),
	IsNew:        false,
}
var ExampleGachaResult2 = &GachaResult{
	CollectionID: "1002",
	ItemName:     "example2",
	Rarity:       int32(2),
	IsNew:        true,
}
var ExampleGachaResult3 = &GachaResult{
	CollectionID: "1003",
	ItemName:     "example3",
	Rarity:       int32(3),
	IsNew:        true,
}

// ExampleGachaProbSlice GachaProb のスライスの例
var ExampleGachaProbSlice = []*gm.GachaProb{
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
