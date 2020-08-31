// Package usecase userシステムのユースケースを満たす処理の流れを実装
package usecase

import (
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	SelectUserLByuserID(string) (*model.UserL, error)
	InsertUserL(*model.UserL) error
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
func (uu userUseCase) SelectUserLByuserID(id string) (user *model.UserL, err error) {
	// persistenceを呼び出す
	user, err = uu.userRepository.SelectUserLByuserID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// InsertUserL Userデータを登録するためのユースケース
func (uu userUseCase) InsertUserL(record *model.UserL) (err error) {
	if err = uu.userRepository.InsertUserL(record); err != nil {
		return err
	}
	return nil
}
