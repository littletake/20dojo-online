package gacha

// import "20dojo-online/pkg/server/domain/model"

// // ChangeToGachaResult 変換する
// func ChangeToGachaResult(gettingItemSlice []string, cItemSlice []*model.CollectionItem,
// 	hasGotItemMap map[string](bool), gachaResultSlice []*model.GachaResult, userID string) []*model.UserCollectionItem {
// 	var newItemSlice []*model.UserCollectionItem

// 	for i, gettingItem := range gettingItemSlice {
// 		for _, item := range cItemSlice {
// 			if gettingItem == item.ItemID {
// 				// 既出itemの確認
// 				if hasGotItemMap[item.ItemID] {
// 					result := model.GachaResult{
// 						CollectionID: item.ItemID,
// 						ItemName:     item.ItemName,
// 						Rarity:       item.Rarity,
// 						IsNew:        false,
// 					}
// 					gachaResultSlice[i] = &result
// 				} else {
// 					result := model.GachaResult{
// 						CollectionID: item.ItemID,
// 						ItemName:     item.ItemName,
// 						Rarity:       item.Rarity,
// 						IsNew:        true,
// 					}
// 					gachaResultSlice[i] = &result
// 					// 既出アイテム一覧に追加
// 					hasGotItemMap[item.ItemID] = true
// 					// 登録
// 					newItem := model.UserCollectionItem{
// 						UserID:           userID,
// 						CollectionItemID: item.ItemID,
// 					}
// 					newItemSlice = append(newItemSlice, &newItem)
// 				}
// 			}
// 		}
// 	}
// 	return newItemSlice
// }
