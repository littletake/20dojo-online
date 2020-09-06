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

// UserHandler UserにおけるHandlerのインターフェース
type UserHandler interface {
	HandleUserGet() http.HandlerFunc
	HandleUserCreate() http.HandlerFunc
	HandleUserUpdate() http.HandlerFunc
}

// userHandler usecaseとhandlerをつなぐもの
type userHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler Handlerを生成する関数
func NewUserHandler(uu usecase.UserUseCase) UserHandler {
	return &userHandler{
		userUseCase: uu,
	}
}

// HandleUserGet ユーザ情報取得
func (uh userHandler) HandleUserGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
			myErr := myerror.NewMyErr(
				fmt.Errorf("userID is empty"),
				500,
			)
			myErr.HandleErr(writer)
			return
		}
		// ユーザ取得
		user, myErr := uh.userUseCase.GetUserByUserID(userID)
		if myErr != nil {
			myErr.HandleErr(writer)
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
}

// HandleUserCreate　ユーザ作成
func (uh userHandler) HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
			myErr := myerror.NewMyErr(err, 500)
			myErr.HandleErr(writer)
			return
		}
		// ユーザを登録
		authToken, myErr := uh.userUseCase.RegisterUserFromUserName(requestBody.Name)
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		// 生成した認証トークンを返却
		response.Success(writer, &userCreateResponse{
			Token: authToken,
		})
	}
}

// HandleUserUpdate　ユーザ情報更新
func (uh userHandler) HandleUserUpdate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// userUpdateRequest ユーザ更新request
		type userUpdateRequest struct {
			Name string `json:"name"`
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
		// requestBodyから更新情報を取得
		var requestBody userUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(err, 500)
			myErr.HandleErr(writer)
			return
		}
		// userNameを更新
		if myErr := uh.userUseCase.UpdateUserName(userID, requestBody.Name); myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		response.Success(writer, nil)
	}
}
