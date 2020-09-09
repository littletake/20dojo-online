package usercollectionitem

import (
	"database/sql"

	model "20dojo-online/pkg/server/domain/model/usercollectionitem"
)

// UCItemRepository UserCollectionItem におけるRepository のインターフェース
type UCItemRepository interface {
	SelectUCItemSliceByUserID(userID string) ([]*model.UserCollectionItem, error)
	BulkInsertUCItemSlice([]*model.UserCollectionItem, *sql.Tx) error
}
