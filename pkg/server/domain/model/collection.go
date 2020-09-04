package model

// CollectionItem コレクションアイテム一覧
type CollectionItemResult struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}
