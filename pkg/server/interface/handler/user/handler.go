package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/user"
)

// UserHandler UserにおけるHandlerのインターフェース
type UserHandler interface {
	HandleUserGet() middleware.MyHandlerFunc
	HandleUserCreate() middleware.MyHandlerFunc
	HandleUserUpdate() middleware.MyHandlerFunc
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

// UserCreateRequest ユーザ作成request
type UserCreateRequest struct {
	Name string `json:"name"`
}

// UserCreateResponse ユーザ作成response
type UserCreateResponse struct {
	Token string `json:"token"`
}

// HandleUserGet ユーザ情報取得
func (uh *userHandler) HandleUserGet() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
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
				http.StatusInternalServerError,
			)
			return myErr
		}
		// ユーザ取得
		user, myErr := uh.userUseCase.GetUserByUserID(userID)
		if myErr != nil {
			return myErr
		}
		// レスポンス
		response.Success(writer, &userGetResponse{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
			Coin:      user.Coin,
		})
		return nil
	}
}

// HandleUserCreate　ユーザ作成
func (uh *userHandler) HandleUserCreate() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
		// リクエストBodyから更新情報を取得
		var requestBody UserCreateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}

		authToken, myErr := uh.userUseCase.RegisterUserFromUserName(requestBody.Name)
		if myErr != nil {
			return myErr
		}
		// 生成した認証トークンを返却
		response.Success(writer, &UserCreateResponse{
			Token: authToken,
		})
		return nil
	}
}

// UserUpdateRequest ユーザ更新request
type UserUpdateRequest struct {
	Name string `json:"name"`
}

// HandleUserUpdate　ユーザ情報更新
func (uh *userHandler) HandleUserUpdate() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {

		// コンテキストからuserID取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("userID is empty"),
				http.StatusInternalServerError,
			)
			return myErr
		}
		// requestBodyから更新情報を取得
		var requestBody UserUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}
		// userNameを更新
		_, myErr := uh.userUseCase.UpdateUserName(userID, requestBody.Name)
		if myErr != nil {
			return myErr
		}
		response.Success(writer, nil)
		return nil
	}
}
