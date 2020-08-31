package repository

import "20dojo-online/pkg/server/domain/model"

// TODO: 命名を考える

// UserRepository User におけるRepository のインターフェース
type UserRepository interface {
	SelectUserLByuserID(string) (*model.UserL, error)
	InsertUserL(*model.UserL) error
}
