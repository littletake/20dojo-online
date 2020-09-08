package persistence

import (
	"database/sql"

	"20dojo-online/pkg/db"
	"20dojo-online/pkg/server/domain/model"
	repository "20dojo-online/pkg/server/domain/repository/collection_item"
)

type cItemPersistence struct{}

// NewCItemPersistence User データに関するPersistence を生成
func NewCItemPersistence() repository.CItemRepository {
	return &cItemPersistence{}
}

func (cp cItemPersistence) SelectAllCollectionItem() ([]*model.CollectionItem, error) {
	rows, err := db.Conn.Query("SELECT * FROM collection_item")
	if err != nil {
		return nil, err
	}
	return convertToItemSlice(rows)

}

// convertToItemSlice rowsデータをCollectionItemSliceデータへ変換する
func convertToItemSlice(rows *sql.Rows) ([]*model.CollectionItem, error) {
	collectionItemSlice := []*model.CollectionItem{}
	for rows.Next() {
		collectionItem := model.CollectionItem{}
		err := rows.Scan(&collectionItem.ItemID, &collectionItem.ItemName, &collectionItem.Rarity)
		if err != nil {
			return nil, err
		}
		collectionItemSlice = append(collectionItemSlice, &collectionItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collectionItemSlice, nil
}
