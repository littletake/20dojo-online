package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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

// TODO: dcontext.GetUserIDFromContextと重複?

// GetUser userの取得
func GetUser(writer http.ResponseWriter, request *http.Request, uh userHandler) (user *model.UserL) {
	// userID取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	// ユーザデータの取得処理と存在チェックを実装
	user, err := uh.userUseCase.SelectUserLByUserID(userID)
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
	return user
}

// CreateUser ユーザ作成
func CreateUser(userName string) (user model.UserL, err error) {
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	// UUIDで認証トークンを生成する
	authToken, err := uuid.NewRandom()
	user = model.UserL{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      userName,
		HighScore: 0,
		Coin:      0,
	}
	return user, err
}

// HandleUserLGet ユーザ情報取得のHandler
func (uh userHandler) HandleUserLGet(writer http.ResponseWriter, request *http.Request) {
	// userGetResponse ユーザ取得response
	type userGetResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		HighScore int32  `json:"highScore"`
		Coin      int32  `json:"coin"`
	}
	// user取得
	user := GetUser(writer, request, uh)
	// レスポンス
	response.Success(writer, &userGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})

}

// HandleUserLCreate
func (uh userHandler) HandleUserLCreate(writer http.ResponseWriter, request *http.Request) {
	// userCreateRequest ユーザ作成request
	type userCreateRequest struct {
		Name string `json:"name"`
	}
	// userCreateResponse ユーザ作成response
	type userCreateResponse struct {
		Token string `json:"token"`
	}

	// リクエストBodyから更新後情報を取得
	var requestBody userCreateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		log.Println(err)
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	// データベースにユーザデータを登録する
	user, err := CreateUser(requestBody.Name)
	if err != nil {
		log.Println(err)
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	if err = uh.userUseCase.InsertUserL(&user); err != nil {
		log.Println(err)
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	// 生成した認証トークンを返却
	response.Success(writer, &userCreateResponse{
		Token: user.AuthToken,
	})

}

// HandleUserLUpdate
func (uh userHandler) HandleUserLUpdate(writer http.ResponseWriter, request *http.Request) {
	// userUpdateRequest ユーザ更新request
	type userUpdateRequest struct {
		Name string `json:"name"`
	}
	// requestBodyから更新情報を取得
	var requestBody userUpdateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		log.Println(err)
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	// user取得
	user := GetUser(writer, request, uh)
	// user情報の更新
	user.Name = requestBody.Name

	// user更新
	if err := uh.userUseCase.UpdateUserLByUser(user); err != nil {
		log.Println(err)
		response.InternalServerError(writer, "Internal Server Error")
		return
	}
	response.Success(writer, nil)
}
