package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type gameFinishRequest struct {
	Score int32 `json:"score"`
}
type gameFinishResponse struct {
	Coin int32 `json:"coin"`
}

// TODO: もう少し綺麗に描けると思う．
// - エラー処理
// - 構造体の定義場所

// HandleGameFinish インゲーム終了処理
func HandleGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody gameFinishRequest
		// Decodeの引数が&である理由: interfaceを引数にしているため
		// error_handling: if err := ~; err != nil {}と書く
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// ユーザデータの取得処理と存在チェックを実装
		user, err := model.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}

		// 送られたscoreからcoinへ変換
		// TODO: coinの計算式を考える
		// 現状: coin = score
		tempCoin := requestBody.Score
		if tempCoin < 0 {
			log.Println(errors.New("score must be positive"))
			response.BadRequest(writer, fmt.Sprintf("score must be positive. score=%d", tempCoin))
			return
		}

		// 所持コインの計算
		user.Coin += tempCoin
		// ハイスコアの処理
		if user.HighScore < requestBody.Score {
			user.HighScore = requestBody.Score
		}
		// DBのtable更新
		// 注意: UpdateUserByPrimaryKeyの引数がuserである点
		if err := model.UpdateUserByPrimaryKey(user); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// レスポンスに必要な情報を詰めて返却
		response.Success(writer, &gameFinishResponse{
			Coin: tempCoin,
		})
	}
}
