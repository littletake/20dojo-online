// Package usecase userシステムのユースケースを満たす処理の流れを実装
package usecase

import (
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	SelectUserByAuth(string) (model.UserD, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

// TODO: 返し方を考える

// NewUserUseCase Userデータに関するUseCaseを生成
func NewUserUseCase(ur repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: ur,
	}
}

// SelectUserByAuth Userデータを条件抽出するためのユースケース
// ここで具体的に実装している？？
func (uu userUseCase) SelectUserByAuth(auth string) (user model.UserD, err error) {
	// persistenceを呼び出す
	user, err = uu.userRepository.SelectUserByAuth(auth)
	if err != nil {
		return nil, err
	}
	return user, nil
}
