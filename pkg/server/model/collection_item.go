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

// CollectionItemList CollectionItemのslice
type CollectionItemList []*CollectionItem

// SelectCollectionItemByItemID IDを条件にレコードを取得する
func SelectCollectionItemByItemID(itemID string) (*CollectionItem, error) {
	// idを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM collection_item WHERE id = ?", itemID)
	return convertToCollectionItem(row)
}

// SelectAllCollectionItem table:collection_itemの全件取得
func SelectAllCollectionItem() (CollectionItemList, error) {
	rows, err := db.Conn.Query("SELECT * FROM collection_item")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToCollectionItemList(rows)
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

// convertToCollectionItemList rowsデータをCollectionItemListデータへ変換する
func convertToCollectionItemList(rows *sql.Rows) (CollectionItemList, error) {
	collectionItemList := CollectionItemList{}
	for rows.Next() {
		collectionItem := CollectionItem{}
		err := rows.Scan(&collectionItem.ItemID, &collectionItem.ItemName, &collectionItem.Rarity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		collectionItemList = append(collectionItemList, &collectionItem)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return collectionItemList, nil
}
