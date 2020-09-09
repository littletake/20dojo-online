package collectionitem

import (
	"database/sql"

	model "20dojo-online/pkg/server/domain/model/collectionitem"
	repository "20dojo-online/pkg/server/domain/repository/collectionitem"
)

type cItemPersistence struct {
	db *sql.DB
}

// NewCItemPersistence User データに関するPersistence を生成
func NewCItemPersistence(db *sql.DB) repository.CItemRepository {
	return &cItemPersistence{
		db: db,
	}
}

func (cp cItemPersistence) SelectAllCollectionItem() ([]*model.CollectionItem, error) {
	rows, err := cp.db.Query("SELECT * FROM collection_item")
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
