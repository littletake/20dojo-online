package repository

import "20dojo-online/pkg/server/domain/model"

// TODO: 命名を考える

// UserRepository User におけるRepository のインターフェース
type UserRepository interface {
	SelectUserByUserID(userID string) (user *model.UserL, err error)
	InsertUser(user *model.UserL) (err error)
	UpdateUserByUser(user *model.UserL) (err error)
	SelectUsersByHighScore(start int32) (user []*model.UserL, err error)
}
