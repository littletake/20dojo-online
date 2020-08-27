package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

// CollectionItem Gachaアイテムの情報
type CollectionItem struct {
	ItemID   string
	ItemName string
	Rarity   int32
}

// SelectCollectionItemByItemID IDを条件にレコードを取得する
func SelectCollectionItemByItemID(itemID string) (*CollectionItem, error) {
	// idを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM collection_item WHERE id = ?", itemID)
	return convertToCollectionItem(row)
}

// convertToCollectionItem rowデータをGachaItemデータへ変換する
func convertToCollectionItem(row *sql.Row) (*CollectionItem, error) {
	collectionItem := CollectionItem{}
	err := row.Scan(&collectionItem.ItemID, &collectionItem.ItemName, &collectionItem.Rarity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &collectionItem, nil
}
