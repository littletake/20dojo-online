//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package user

import (
	"database/sql"

	model "20dojo-online/pkg/server/domain/model/user"
)

// UserRepo User におけるRepository のインターフェース
type UserRepo interface {
	SelectUserByUserID(userID string) (*model.UserL, error)
	SelectUserByAuthToken(userID string) (*model.UserL, error)
	SelectUsersByHighScore(limit int32, start int32) ([]*model.UserL, error)
	InsertUser(user *model.UserL) error
	UpdateUserByUser(user *model.UserL) error
	UpdateUserByUserInTx(user *model.UserL, tx *sql.Tx) error
}
