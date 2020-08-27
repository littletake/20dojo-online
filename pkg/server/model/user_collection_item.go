package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

// UserCollectionItem table: user_collection_itemのデータ
type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
}

// InsertCollectionItem データベースにレコードを登録する
func InsertCollectionItem(record *UserCollectionItem) error {
	// userテーブルへのレコードの登録を行うSQLを入力する
	stmt, err := db.Conn.Prepare("INSERT INTO user_collection_item (id, collection_item_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.UserID, record.CollectionItemID)
	return err
}

// TODO: 命名を考える

// SelectUserCollectionItemByItemIDAndUserID itemIDとuserIDを条件にレコードを取得する
func SelectUserCollectionItemByItemIDAndUserID(itemID string, userID string) (*UserCollectionItem, error) {
	row := db.Conn.QueryRow("SELECT * FROM user_collection_item WHERE user_id = ? AND collection_item_id = ?", userID, itemID)
	// TODO: 参照するアイテムがないのかerrなのかを判定できるようにする
	return convertToUserCollectionItem(row)
}

// convertToUserCollectionItem rowデータをUserCollectionItemデータへ変換する
func convertToUserCollectionItem(row *sql.Row) (*UserCollectionItem, error) {
	userCollectionItem := UserCollectionItem{}
	err := row.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &userCollectionItem, nil
}
