package repository

import (
	"database/sql"

	"20dojo-online/pkg/server/domain/model"
)

// UserRepository User におけるRepository のインターフェース
type UserRepository interface {
	SelectUserByUserID(userID string) (*model.UserL, error)
	SelectUserByAuthToken(userID string) (*model.UserL, error)
	SelectUsersByHighScore(start int32) ([]*model.UserL, error)
	InsertUser(user *model.UserL) error
	UpdateUserByUser(user *model.UserL) error
	UpdateUserByUserInTx(user *model.UserL, tx *sql.Tx) error
}
