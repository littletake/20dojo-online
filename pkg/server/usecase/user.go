// Package usecase userシステムのユースケースを満たす処理の流れを実装
package usecase

import (
	"fmt"

	"github.com/google/uuid"

	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	GetUserLByUserID(string) (*model.UserL, *model.MyErr)
	RegisterUser(string) (string, *model.MyErr)
	UpdateUserLByUser(*model.UserL) error
	CreateMyErr(error, int32) *model.MyErr
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

// CreateMyErr MyErrの定義
func (uu userUseCase) CreateMyErr(errMsg error, errCode int32) *model.MyErr {
	myErr := &model.MyErr{
		ErrMsg:  errMsg,
		ErrCode: errCode,
	}
	return myErr
}

// GetUserByAuth Userデータを条件抽出するためのユースケース
func (uu userUseCase) GetUserLByUserID(id string) (user *model.UserL, myErr *model.MyErr) {
	// idと照合するユーザを取得
	user, err := uu.userRepository.SelectUserLByUserID(id)
	if err != nil {
		// TODO: こうやって使っていいのか？
		myErr = uu.CreateMyErr(err, 500)
		return nil, myErr
	}
	if user == nil {
		myErr = uu.CreateMyErr(
			fmt.Errorf("user not found"),
			500,
		)
		return nil, myErr
	}
	return user, nil
}

// RegisterUser Userデータを登録するためのユースケース
func (uu userUseCase) RegisterUser(userName string) (string, *model.MyErr) {
	// TODO: どのIDの生成でエラーが生じたのかをエラーメッセージに添付すること
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		myErr := uu.CreateMyErr(err, 500)
		return "", myErr
	}
	// UUIDで認証トークンを生成する
	authToken, err := uuid.NewRandom()
	if err != nil {
		myErr := uu.CreateMyErr(err, 500)
		return "", myErr
	}
	// ユーザ作成
	user := &model.UserL{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      userName,
		HighScore: 0,
		Coin:      0,
	}
	// ユーザ登録
	if err = uu.userRepository.InsertUserL(user); err != nil {
		myErr := uu.CreateMyErr(err, 500)
		return "", myErr
	}
	return user.AuthToken, nil
}

// UpdateUserLByUser Userデータを更新するためのユースケース
func (uu userUseCase) UpdateUserLByUser(record *model.UserL) (err error) {
	if err = uu.userRepository.UpdateUserLByUser(record); err != nil {
		return err
	}
	return nil
}
