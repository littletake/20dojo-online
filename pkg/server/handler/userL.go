package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/usecase"
)

// UserHandler UserにおけるHandlerのインターフェース
type UserHandler interface {
	HandleUserLGet(http.ResponseWriter, *http.Request)
	HandleUserLCreate(http.ResponseWriter, *http.Request)
	HandleUserLUpdate(http.ResponseWriter, *http.Request)
	HandleGameFinish(http.ResponseWriter, *http.Request)
}

// TODO: あまりわかっていない
// userHandler usecaseとhandlerをつなぐもの
type userHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler Userデータに関するHandler
func NewUserHandler(uu usecase.UserUseCase) UserHandler {
	return &userHandler{
		userUseCase: uu,
	}
}

// DoErr エラー時の処理
func DoErr(writer http.ResponseWriter, myErr *model.MyErr) {
	if myErr.ErrCode == 400 {
		log.Println(myErr.ErrMsg)
		response.BadRequest(writer, myErr.ErrMsg.Error())
		return
	} else if myErr.ErrCode == 500 {
		log.Println(myErr.ErrMsg)
		response.InternalServerError(writer, myErr.ErrMsg.Error())
		return
	} else {
		// TODO: エラーコードが400,500以外の場合の処理考える
		errMsg := fmt.Sprintf("!! errorCode mistake. ErrCode: %d !!", myErr.ErrCode)
		panic(errMsg)
	}

}

// HandleUserLGet ユーザ情報取得
func (uh userHandler) HandleUserLGet(writer http.ResponseWriter, request *http.Request) {
	// userGetResponse ユーザ取得response
	type userGetResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		HighScore int32  `json:"highScore"`
		Coin      int32  `json:"coin"`
	}

	// コンテキストからuserID取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		myErr := uh.userUseCase.CreateMyErr(
			fmt.Errorf("userID is empty"),
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// ユーザ取得
	user, myErr := uh.userUseCase.GetUserLByUserID(userID)
	if myErr != nil {
		DoErr(writer, myErr)
		return
	}
	// レスポンス
	response.Success(writer, &userGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})

}

// HandleUserLCreate　ユーザ作成
func (uh userHandler) HandleUserLCreate(writer http.ResponseWriter, request *http.Request) {
	// userCreateRequest ユーザ作成request
	type userCreateRequest struct {
		Name string `json:"name"`
	}
	// userCreateResponse ユーザ作成response
	type userCreateResponse struct {
		Token string `json:"token"`
	}

	// リクエストBodyから更新情報を取得
	var requestBody userCreateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		myErr := uh.userUseCase.CreateMyErr(
			err,
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// ユーザを登録
	authToken, myErr := uh.userUseCase.RegisterUserFromUserName(requestBody.Name)
	if myErr != nil {
		DoErr(writer, myErr)
		return
	}
	// 生成した認証トークンを返却
	response.Success(writer, &userCreateResponse{
		Token: authToken,
	})

}

// HandleUserLUpdate　ユーザ情報更新
func (uh userHandler) HandleUserLUpdate(writer http.ResponseWriter, request *http.Request) {
	// userUpdateRequest ユーザ更新request
	type userUpdateRequest struct {
		Name string `json:"name"`
	}

	// コンテキストからuserID取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		myErr := uh.userUseCase.CreateMyErr(
			fmt.Errorf("userID is empty"),
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// requestBodyから更新情報を取得
	var requestBody userUpdateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		myErr := uh.userUseCase.CreateMyErr(
			err,
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// userNameを更新
	myErr := uh.userUseCase.UpdateUserName(userID, requestBody.Name)
	if myErr != nil {
		DoErr(writer, myErr)
		return
	}
	response.Success(writer, nil)
}

// HandleGameFinish インゲーム終了処理
func (uh userHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {
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
		myErr := uh.userUseCase.CreateMyErr(
			fmt.Errorf("userID is empty"),
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// リクエストBodyから更新後情報を取得
	var requestBody GameFinishRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		myErr := uh.userUseCase.CreateMyErr(
			err,
			500,
		)
		DoErr(writer, myErr)
		return
	}
	// coinとscoreを更新
	coin, myErr := uh.userUseCase.UpdateCoinAndHighScore(userID, requestBody.Score)
	if myErr != nil {
		DoErr(writer, myErr)
		return
	}
	// レスポンスに必要な情報を詰めて返却
	response.Success(writer, &GameFinishResponse{
		Coin: coin,
	})
}
