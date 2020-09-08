package user_collection_item

import (
	"database/sql"

	"20dojo-online/pkg/server/domain/model"
)

// UCItemRepository UserCollectionItem におけるRepository のインターフェース
type UCItemRepository interface {
	SelectUCItemSliceByUserID(userID string) ([]*model.UserCollectionItem, error)
	BulkInsertUCItemSlice([]*model.UserCollectionItem, *sql.Tx) error
}
