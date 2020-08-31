package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

// CollectionItem table:collection_itemの内容
type CollectionItem struct {
	ItemID   string
	ItemName string
	Rarity   int32
}

// CollectionItemSlice CollectionItemのslice
type CollectionItemSlice []*CollectionItem

// SelectAllCollectionItem table:collection_itemの全件取得
func SelectAllCollectionItem() (CollectionItemSlice, error) {
	rows, err := db.Conn.Query("SELECT * FROM collection_item")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToCollectionItemSlice(rows)
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

// convertToCollectionItemSlice rowsデータをCollectionItemSliceデータへ変換する
func convertToCollectionItemSlice(rows *sql.Rows) (CollectionItemSlice, error) {
	collectionItemSlice := CollectionItemSlice{}
	for rows.Next() {
		collectionItem := CollectionItem{}
		err := rows.Scan(&collectionItem.ItemID, &collectionItem.ItemName, &collectionItem.Rarity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		collectionItemSlice = append(collectionItemSlice, &collectionItem)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return collectionItemSlice, nil
}
