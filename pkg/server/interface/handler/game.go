package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	"20dojo-online/pkg/server/usecase"
)

// GameHandler gameにおけるHandler
type GameHandler interface {
	HandleGameFinish() http.HandlerFunc
}

// gameHandler usecaseとhandlerをつなぐもの
type gameHandler struct {
	gameUseCase usecase.GameUseCase
}

// NewGameHandler Handlerを生成
func NewGameHandler(gu usecase.GameUseCase) GameHandler {
	return &gameHandler{
		gameUseCase: gu,
	}
}

// HandleGameFinish インゲーム終了処理
func (gh gameHandler) HandleGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// GameFinishRequest ゲーム終了request
		type GameFinishRequest struct {
			Score int32 `json:"score"`
		}
		// GameFinishResponse ゲーム終了response
		type GameFinishResponse struct {
			Coin int32 `json:"coin"`
		}

		// コンテキストからuserID取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("userID is empty"),
				500,
			)
			myErr.HandleErr(writer)
			return
		}
		// リクエストBodyから更新後情報を取得
		var requestBody GameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(err, 500)
			myErr.HandleErr(writer)
			return
		}
		// coinとscoreを更新
		coin, myErr := gh.gameUseCase.UpdateCoinAndHighScore(userID, requestBody.Score)
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		// レスポンスに必要な情報を詰めて返却
		response.Success(writer, &GameFinishResponse{
			Coin: coin,
		})
	}
}
