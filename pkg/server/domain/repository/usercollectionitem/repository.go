package usercollectionitem

import (
	"database/sql"

	model "20dojo-online/pkg/server/domain/model/usercollectionitem"
)

// UserCollectionItemRepo UserCollectionItem におけるRepository のインターフェース
type UserCollectionItemRepo interface {
	SelectSliceByUserID(userID string) ([]*model.UserCollectionItem, error)
	BulkInsert([]*model.UserCollectionItem, *sql.Tx) error
}
