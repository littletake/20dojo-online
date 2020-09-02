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
	CreateMyErr(errMsg error, errCode int32) (myErr *model.MyErr)
	GetUserLByUserID(userID string) (user *model.UserL, myErr *model.MyErr)
	RegisterUserFromUserName(userName string) (authToken string, myErr *model.MyErr)
	UpdateUserName(userID string, userName string) (myErr *model.MyErr)
	UpdateCoinAndHighScore(userID string, score int32) (coin int32, myErr *model.MyErr)
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

// CreateMyErr MyErrの定義
func (uu userUseCase) CreateMyErr(errMsg error, errCode int32) (myErr *model.MyErr) {
	myErr = &model.MyErr{
		ErrMsg:  errMsg,
		ErrCode: errCode,
	}
	return myErr
}

// GetUserByAuth Userデータを条件抽出
func (uu userUseCase) GetUserLByUserID(userID string) (user *model.UserL, myErr *model.MyErr) {
	// idと照合するユーザを取得
	user, err := uu.userRepository.SelectUserLByUserID(userID)
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

// RegisterUserFromUserName Userデータを登録
func (uu userUseCase) RegisterUserFromUserName(userName string) (authToken string, myErr *model.MyErr) {
	// TODO: どのIDの生成でエラーが生じたのかをエラーメッセージに添付すること
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		myErr := uu.CreateMyErr(err, 500)
		return "", myErr
	}
	// UUIDで認証トークンを生成する
	token, err := uuid.NewRandom()
	if err != nil {
		myErr := uu.CreateMyErr(err, 500)
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
	if err = uu.userRepository.InsertUserL(user); err != nil {
		myErr = uu.CreateMyErr(err, 500)
		return "", myErr
	}
	return user.AuthToken, nil
}

// UpdateUserName UserNameを更新
func (uu userUseCase) UpdateUserName(userID string, userName string) (myErr *model.MyErr) {
	// ユーザ取得
	user, myErr := uu.GetUserLByUserID(userID)
	if myErr != nil {
		return myErr
	}
	// ユーザ更新
	user.Name = userName
	// 更新を保存
	if err := uu.userRepository.UpdateUserLByUser(user); err != nil {
		myErr := uu.CreateMyErr(
			err,
			500,
		)
		return myErr
	}
	return nil
}

// UpdateCoinAndHighScore CoinとScoreを更新
func (uu userUseCase) UpdateCoinAndHighScore(userID string, score int32) (coin int32, myErr *model.MyErr) {
	// coinとhighScoreを更新
	if score < 0 {
		myErr := uu.CreateMyErr(
			fmt.Errorf("score must be positive. score=%d", score),
			400,
		)
		return 0, myErr
	}
	// ユーザ取得
	user, myErr := uu.GetUserLByUserID(userID)
	if myErr != nil {
		return 0, myErr
	}
	// コインに変換
	coin = ChangeScoreToCoin(score)
	// 所持コインの追加
	user.Coin += coin
	// ハイスコアの処理
	if user.HighScore < score {
		user.HighScore = score
	}
	// 更新を保存
	if err := uu.userRepository.UpdateUserLByUser(user); err != nil {
		myErr := uu.CreateMyErr(
			err,
			500,
		)
		return 0, myErr
	}
	return coin, nil
}
