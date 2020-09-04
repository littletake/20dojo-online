package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/usecase"
)

// GameHandler gameにおけるHandler
type GameHandler interface {
	HandleGameFinish(writer http.ResponseWriter, request *http.Request)
}

// gameHandler usecaseとhandlerをつなぐもの
type gameHandler struct {
	gameUseCase usecase.GameUseCase
}

// NewGameHandler Userデータに関するHandler
func NewGameHandler(gu usecase.GameUseCase) GameHandler {
	return &gameHandler{
		gameUseCase: gu,
	}
}

// HandleGameFinish インゲーム終了処理
func (gh gameHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {
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
		myErr := myerror.MyErr{
			fmt.Errorf("userID is empty"),
			500,
		}
		myErr.HandleErr(writer)
		return
	}
	// リクエストBodyから更新後情報を取得
	var requestBody GameFinishRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		myErr := myerror.MyErr{err, 500}
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

// // HandleGameFinish インゲーム終了処理
// func (uh userHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {
// 	// GameFinishRequest ゲーム終了request
// 	type GameFinishRequest struct {
// 		Score int32 `json:"score"`
// 	}
// 	// GameFinishResponse ゲーム終了response
// 	type GameFinishResponse struct {
// 		Coin int32 `json:"coin"`
// 	}

// 	// コンテキストからuserID取得
// 	ctx := request.Context()
// 	userID := dcontext.GetUserIDFromContext(ctx)
// 	if userID == "" {
// 		myErr := useCase.CreateMyErr(
// 			fmt.Errorf("userID is empty"),
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// リクエストBodyから更新後情報を取得
// 	var requestBody GameFinishRequest
// 	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
// 		myErr := useCase.CreateMyErr(
// 			err,
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// coinとscoreを更新
// 	coin, myErr := useCase.UpdateCoinAndHighScore(userID, requestBody.Score)
// 	if myErr != nil {
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// レスポンスに必要な情報を詰めて返却
// 	response.Success(writer, &GameFinishResponse{
// 		Coin: coin,
// 	})
// }
