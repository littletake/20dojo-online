package persistence

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"20dojo-online/pkg/db"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

type ucItemPersistence struct{}

// NewUCItemPersistence UserCollectionItem データに関するPersistence を生成
func NewUCItemPersistence() repository.UCItemRepository {
	return &ucItemPersistence{}
}

// SelectUCItemSliceByUserID userIDを条件にレコードを取得する
func (up ucItemPersistence) SelectUCItemSliceByUserID(userID string) ([]*model.UserCollectionItem, error) {
	rows, err := db.Conn.Query("SELECT * FROM user_collection_item WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	return convertToUCItemSlice(rows)
}

// TODO: トランザクションを使わない場合も考慮して実装すること

// BulkInsertUserCollectionItem データベースに複数レコードを登録する
func (up ucItemPersistence) BulkInsertUCItemSlice(records []*model.UserCollectionItem, tx *sql.Tx) error {
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
	stmt, err := tx.Prepare(queryString.String())
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// convertToUCItemSlice rowsデータをUserCollectionItemSliceデータへ変換する
func convertToUCItemSlice(rows *sql.Rows) ([]*model.UserCollectionItem, error) {
	userCollectionItemSlice := []*model.UserCollectionItem{}
	for rows.Next() {
		userCollectionItem := model.UserCollectionItem{}
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
