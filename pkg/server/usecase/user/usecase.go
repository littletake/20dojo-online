//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package user

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	model "20dojo-online/pkg/server/domain/model/user"
	ur "20dojo-online/pkg/server/domain/repository/user"
	"20dojo-online/pkg/server/interface/myerror"
)

// UserUseCase UserにおけるUseCaseのインターフェース
type UserUseCase interface {
	GetUserByUserID(userID string) (*model.UserL, *myerror.MyErr)
	GetUserByAuthToken(token string) (*model.UserL, *myerror.MyErr)
	RegisterUserFromUserName(userName string) (string, *myerror.MyErr)
	UpdateUserName(userID string, userName string) (*model.UserL, *myerror.MyErr)
}

type userUseCase struct {
	userRepository ur.UserRepo
	createUUID     func() (uuid.UUID, error)
}

// NewUserUseCase Userデータに関するUseCaseを生成
func NewUserUseCase(ur ur.UserRepo, f func() (uuid.UUID, error)) UserUseCase {
	return &userUseCase{
		userRepository: ur,
		createUUID:     f,
	}
}

// GetUserByUserID Userデータを条件抽出
func (uu *userUseCase) GetUserByUserID(userID string) (*model.UserL, *myerror.MyErr) {
	// idと照合するユーザを取得
	user, err := uu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found"),
			500,
		)
		return nil, myErr
	}
	return user, nil
}

// GetUserByAuthToken
func (uu *userUseCase) GetUserByAuthToken(token string) (*model.UserL, *myerror.MyErr) {
	// tokenと照合するユーザを取得
	user, err := uu.userRepository.SelectUserByAuthToken(token)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found. token=%s", token),
			400,
		)
		return nil, myErr
	}
	return user, nil
}

// RegisterUserFromUserName Userデータを登録
func (uu *userUseCase) RegisterUserFromUserName(userName string) (string, *myerror.MyErr) {
	// userID
	userID, err := uu.createUUID()
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return "", myErr
	}
	// token
	token, err := uu.createUUID()
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return "", myErr
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
	if err := uu.userRepository.InsertUser(user); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return "", myErr
	}
	return user.AuthToken, nil
}

// UpdateUserName UserNameを更新
func (uu *userUseCase) UpdateUserName(userID string, userName string) (*model.UserL, *myerror.MyErr) {
	// ユーザ取得
	// idと照合するユーザを取得
	user, err := uu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found"),
			500,
		)
		return nil, myErr
	}
	// ユーザ更新
	user.Name = userName
	// 更新を保存
	if err := uu.userRepository.UpdateUserByUser(user); err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	return user, nil
}
