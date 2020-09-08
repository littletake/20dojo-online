package middleware

import (
	"context"
	"fmt"
	"net/http"

	"20dojo-online/pkg/server/interface/myerror"

	"20dojo-online/pkg/dcontext"
)

// // Middleware middlewareのインターフェース
// type Middleware interface {
// 	Authenticate(http.HandlerFunc) http.HandlerFunc
// }

// type middleware struct {
// 	userUseCase usecase.UserUseCase
// }

// // NewMiddleware userUseCaseと疎通
// func NewMiddleware(uu usecase.UserUseCase) Middleware {
// 	return &middleware{
// 		userUseCase: uu,
// 	}
// }

func (m middleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("x-token is empty"),
				400,
			)
			myErr.HandleErr(writer)
			return
		}
		// データベースから認証トークンに紐づくユーザの情報を取得
		user, myErr := m.userUseCase.GetUserByAuthToken(token)
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 認証を通過して，引数に書いた処理が実行
		nextFunc(writer, request.WithContext(ctx))
	}
}
