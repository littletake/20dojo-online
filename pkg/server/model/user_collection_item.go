package model

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"20dojo-online/pkg/db"
)

// UserCollectionItem table:user_collection_itemのデータ
type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
}

// UserCollectionItemSlice UserCollectionItemのslice
type UserCollectionItemSlice []*UserCollectionItem

// InsertUserCollectionItem データベースにレコードを登録する
func InsertUserCollectionItem(record *UserCollectionItem) error {
	// userテーブルへのレコードの登録を行うSQLを入力する
	stmt, err := db.Conn.Prepare("INSERT INTO user_collection_item (id, collection_item_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.UserID, record.CollectionItemID)
	return err
}

// BulkInsertUserCollectionItem データベースに複数レコードを登録する
func BulkInsertUserCollectionItem(records UserCollectionItemSlice) error {
	var queryString strings.Builder
	queryString.WriteString("INSERT INTO user_collection_item (user_id, collection_item_id) VALUES ")
	for i, record := range records {
		queryString.WriteString("(")
		queryString.WriteString(strconv.Quote(record.UserID))
		queryString.WriteString(", ")
		queryString.WriteString(strconv.Quote(record.CollectionItemID))
		queryString.WriteString(")")
		if i != len(records)-1 {
			queryString.WriteString(", ")
		}
	}
	stmt, err := db.Conn.Prepare(queryString.String())
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// SelectUserCollectionItemByItemIDAndUserID itemIDとuserIDを条件にレコードを取得する
func SelectUserCollectionItemByItemIDAndUserID(itemID string, userID string) (*UserCollectionItem, error) {
	row := db.Conn.QueryRow("SELECT * FROM user_collection_item WHERE user_id = ? AND collection_item_id = ?", userID, itemID)
	// TODO: 参照するアイテムがないのかerrなのかを判定できるようにする
	return convertToUserCollectionItem(row)
}

// SelectAllUserCollectionItem table:user_collection_itemの全件取得
func SelectAllUserCollectionItem() (UserCollectionItemSlice, error) {
	rows, err := db.Conn.Query("SELECT * FROM user_collection_item")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToUserCollectionItemSlice(rows)
}

// SelectUserCollectionItemSliceByUserID userIDを条件にレコードを取得する
func SelectUserCollectionItemSliceByUserID(userID string) (UserCollectionItemSlice, error) {
	rows, err := db.Conn.Query("SELECT * FROM user_collection_item WHERE user_id = ?", userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToUserCollectionItemSlice(rows)
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

// convertToUserCollectionItemSlice rowsデータをUserCollectionItemSliceデータへ変換する
func convertToUserCollectionItemSlice(rows *sql.Rows) (UserCollectionItemSlice, error) {
	userCollectionItemSlice := UserCollectionItemSlice{}
	for rows.Next() {
		userCollectionItem := UserCollectionItem{}
		err := rows.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionItemID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		userCollectionItemSlice = append(userCollectionItemSlice, &userCollectionItem)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return userCollectionItemSlice, nil
}
