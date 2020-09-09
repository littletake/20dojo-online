package collectionitem

import (
	"database/sql"

	"20dojo-online/pkg/db"
	model "20dojo-online/pkg/server/domain/model/collectionitem"
	repository "20dojo-online/pkg/server/domain/repository/collectionitem"
)

type cItemPersistence struct{}

// NewPersistence CollectionItem データに関するPersistence を生成
func NewPersistence() repository.CollectionItemRepo {
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
