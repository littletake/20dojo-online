package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/usecase"
)

// UserHandler UserにおけるHandlerのインターフェース
type UserHandler interface {
	HandleUserLGet(http.ResponseWriter, *http.Request)
}

type userHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler Userデータに関するHandler
func NewUserHandler(uu usecase.UserUseCase) UserHandler {
	return &userHandler{
		userUseCase: uu,
	}
}

// HandleUserLGet
func (uh userHandler) HandleUserLGet(writer http.ResponseWriter, request *http.Request) {
	// UserGetResponse ユーザ取得response
	type UserGetResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		HighScore int32  `json:"highScore"`
		Coin      int32  `json:"coin"`
	}

	// Contextから認証済みのユーザIDを取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		response.InternalServerError(writer, "Internal Server Error")
		return
	}

	// ユースケースの呼び出し
	user, err := uh.userUseCase.SelectUserLByuserID(userID)
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

	// レスポンスに必要な情報を詰めて返却
	response.Success(writer, &UserGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})

}

// // HandleUserLCreate
// func (uh userHandler) HandleUserLCreate() http.HandlerFunc {
// 		// リクエストBodyから更新後情報を取得
// 		var requestBody UserCreateRequest
// 		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
// 			log.Println(err)
// 			response.InternalServerError(writer, "Internal Server Error")
// 		}

// 		// UUIDでユーザIDを生成する
// 		userID, err := uuid.NewRandom()
// 		if err != nil {
// 			log.Println(err)
// 			response.InternalServerError(writer, "Internal Server Error")
// 			return
// 		}

// 		// UUIDで認証トークンを生成する
// 		authToken, err := uuid.NewRandom()
// 		if err != nil {
// 			log.Println(err)
// 			response.InternalServerError(writer, "Internal Server Error")
// 			return
// 		}

// 		// データベースにユーザデータを登録する
// 		// ユーザデータの登録クエリを入力する
// 		// TODO: 既存ユーザが存在していた場合の処理を考えるべき？？（現状，重複登録になる）
// 		err = model.InsertUser(&model.User{
// 			ID:        userID.String(),
// 			AuthToken: authToken.String(),
// 			Name:      requestBody.Name,
// 			HighScore: 0,
// 			Coin:      0,
// 		})
// 		if err != nil {
// 			log.Println(err)
// 			response.InternalServerError(writer, "Internal Server Error")
// 			return
// 		}

// 		// 生成した認証トークンを返却
// 		response.Success(writer, &UserCreateResponse{Token: authToken.String()})
// 	}
// }
