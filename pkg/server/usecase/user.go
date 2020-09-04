// Package usecase userシステムのユースケースを満たす処理の流れを実装
package usecase

import (
	"fmt"

	"github.com/google/uuid"

	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
	"20dojo-online/pkg/server/interface/myerror"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	GetUserByUserID(userID string) (*model.UserL, *myerror.MyErr)
	GetUserByAuthToken(token string) (*model.UserL, *myerror.MyErr)
	RegisterUserFromUserName(userName string) (string, *myerror.MyErr)
	UpdateUserName(userID string, userName string) *myerror.MyErr
}

type userUseCase struct {
	userRepository repository.UserRepository
}

// NewUserUseCase Userデータに関するUseCaseを生成
func NewUserUseCase(ur repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: ur,
	}
}

// GetUserByUserID Userデータを条件抽出
func (uu userUseCase) GetUserByUserID(userID string) (*model.UserL, *myerror.MyErr) {
	// idと照合するユーザを取得
	user, err := uu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return nil, &myErr
	}
	if user == nil {
		myErr := myerror.MyErr{
			fmt.Errorf("user not found"),
			500,
		}
		return nil, &myErr
	}
	return user, nil
}

// GetUserByAuthToken
func (uu userUseCase) GetUserByAuthToken(token string) (*model.UserL, *myerror.MyErr) {
	// tokenと照合するユーザを取得
	user, err := uu.userRepository.SelectUserByAuthToken(token)
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return nil, &myErr
	}
	if user == nil {
		myErr := myerror.MyErr{
			fmt.Errorf("user not found. token=%s", token),
			400,
		}
		return nil, &myErr
	}
	return user, nil
}

// RegisterUserFromUserName Userデータを登録
func (uu userUseCase) RegisterUserFromUserName(userName string) (string, *myerror.MyErr) {
	// TODO: どのIDの生成でエラーが生じたのかをエラーメッセージに添付すること
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return "", &myErr
	}
	// UUIDで認証トークンを生成する
	token, err := uuid.NewRandom()
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return "", &myErr
	}
	// ユーザ作成
	user := &model.UserL{
		ID:        userID.String(),
		AuthToken: token.String(),
		Name:      userName,
		HighScore: 0,
		Coin:      0,
	}
	// ユーザ登録
	if err = uu.userRepository.InsertUser(user); err != nil {
		myErr := myerror.MyErr{err, 500}
		return "", &myErr
	}
	return user.AuthToken, nil
}

// UpdateUserName UserNameを更新
func (uu userUseCase) UpdateUserName(userID string, userName string) *myerror.MyErr {
	// ユーザ取得
	user, myErr := uu.GetUserByUserID(userID)
	if myErr != nil {
		return myErr
	}
	// ユーザ更新
	user.Name = userName
	// 更新を保存
	if err := uu.userRepository.UpdateUserByUser(user); err != nil {
		myErr := myerror.MyErr{err, 500}
		return &myErr
	}
	return nil
}
