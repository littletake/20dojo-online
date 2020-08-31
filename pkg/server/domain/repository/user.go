package repository

import "20dojo-online/pkg/server/domain/model"

// UserRepository User におけるRepository のインターフェース
type UserRepository interface {
	SelectUserByAuth(string) (model.UserD, error)
}
