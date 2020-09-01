// Package usecase userシステムのユースケースを満たす処理の流れを実装
package usecase

import (
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	SelectUserLByUserID(string) (*model.UserL, error)
	InsertUserL(*model.UserL) error
	UpdateUserLByUser(*model.UserL) error
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
func (uu userUseCase) SelectUserLByUserID(id string) (user *model.UserL, err error) {
	// persistenceを呼び出す
	user, err = uu.userRepository.SelectUserLByUserID(id)
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

// UpdateUserLByUser Userデータを更新するためのユースケース
func (uu userUseCase) UpdateUserLByUser(record *model.UserL) (err error) {
	if err = uu.userRepository.UpdateUserLByUser(record); err != nil {
		return err
	}
	return nil
}
