package handler

import (
	"fmt"

	"20dojo-online/pkg/server/domain/model"
)

// // GetUser userの取得
// func GetUser(writer http.ResponseWriter, request *http.Request, uh userHandler) (user *model.UserL, myErr *MyErr) {
// 	// userID取得
// 	ctx := request.Context()
// 	userID := dcontext.GetUserIDFromContext(ctx)
// 	if userID == "" {
// 		myErr = &MyErr{
// 			ErrMsg:  fmt.Errorf("userID is empty"),
// 			ErrCode: 500,
// 		}
// 		return nil, myErr
// 	}
// 	// ユーザデータの取得処理と存在チェックを実装
// 	user, err := uh.userUseCase.SelectUserLByUserID(userID)
// 	if err != nil {
// 		myErr = &MyErr{
// 			ErrMsg:  err,
// 			ErrCode: 500,
// 		}
// 		return nil, myErr
// 	}
// 	if user == nil {
// 		myErr = &MyErr{
// 			ErrMsg:  fmt.Errorf("user not found"),
// 			ErrCode: 500,
// 		}
// 		return nil, myErr
// 	}
// 	return user, nil
// }

// UpdateCoinAndHighScore スコア更新
func UpdateCoinAndHighScore(score int32, user *model.UserL) (coin int32, err string) {
	// TODO: errorをラップして必要な情報をかえす
	if score < 0 {
		errMsg := fmt.Sprintf("score must be positive. score=%d", score)
		return 0, errMsg
	}
	coin = changeScoreToCoin(score)
	// 所持コインの計算
	user.Coin += coin
	// ハイスコアの処理
	if user.HighScore < score {
		user.HighScore = score
	}
	return coin, ""
}

// score -> coin
func changeScoreToCoin(score int32) int32 {
	// 単純に返すだけ
	coin := score
	return coin
}
