//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package game

import (
	"fmt"
	"net/http"

	ur "20dojo-online/pkg/server/domain/repository/user"

	"20dojo-online/pkg/server/interface/myerror"
)

// GameUseCase gameに関するusecase
type GameUseCase interface {
	UpdateCoinAndHighScore(userID string, score int32) (int32, *myerror.MyErr)
}
type gameUseCase struct {
	userRepository ur.UserRepo
}

// NewGameUseCase GameUsecaseの生成
func NewGameUseCase(ur ur.UserRepo) GameUseCase {
	return &gameUseCase{
		userRepository: ur,
	}
}

// ChangeScoreToCoin score から coin に変換
func ChangeScoreToCoin(score int32) int32 {
	// 単純に返すだけ
	coin := score
	return coin
}

// UpdateCoinAndHighScore CoinとScoreを更新
func (gu *gameUseCase) UpdateCoinAndHighScore(userID string, score int32) (int32, *myerror.MyErr) {
	// coinとhighScoreを更新
	if score < 0 {
		myErr := myerror.NewMyErr(
			fmt.Errorf("score must be positive. score=%d", score),
			http.StatusBadRequest,
		)
		return 0, myErr
	}
	// ユーザ取得
	user, err := gu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return 0, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found. userID=%s", userID),
			http.StatusBadRequest,
		)
		return 0, myErr
	}
	// コインに変換
	// TODO: コインに変換するアルゴリズムを工夫する
	coin := ChangeScoreToCoin(score)
	// 所持コインの追加
	user.Coin += coin
	// ハイスコアの処理
	if user.HighScore < score {
		user.HighScore = score
	}
	// 更新を保存
	if err := gu.userRepository.UpdateUserByUser(user); err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return 0, myErr
	}
	return coin, nil
}
