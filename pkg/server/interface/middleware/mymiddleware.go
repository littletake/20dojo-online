package middleware

import (
	"context"
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/myerror"
	usecase "20dojo-online/pkg/server/usecase/user"
)

// MyMiddleware middlewareのインターフェース
type MyMiddleware interface {
	Get(MyHandlerFunc) http.HandlerFunc
	Post(MyHandlerFunc) http.HandlerFunc
	Authenticate(MyHandlerFunc) MyHandlerFunc
}

type myMiddleware struct {
	userUseCase usecase.UserUseCase
}

// NewMyMiddleware userUseCaseと疎通
func NewMyMiddleware(uu usecase.UserUseCase) MyMiddleware {
	return &myMiddleware{
		userUseCase: uu,
	}
}

// Authenticate 認証
func (m *myMiddleware) Authenticate(nextFunc MyHandlerFunc) MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("x-token is empty"),
				http.StatusBadRequest,
			)
			return myErr
		}
		// データベースから認証トークンに紐づくユーザの情報を取得
		user, myErr := m.userUseCase.GetUserByAuthToken(token)
		if myErr != nil {
			return myErr
		}
		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 認証を通過して，引数に書いた処理が実行
		if myErr := nextFunc(writer, request.WithContext(ctx)); myErr != nil {
			return myErr
		}
		return nil
	}
}

// Get GETリクエストを処理する
func (m *myMiddleware) Get(apiFunc MyHandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// Post POSTリクエストを処理する
func (m *myMiddleware) Post(apiFunc MyHandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc MyHandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")

		// エラーハンドリング
		if myErr := apiFunc(writer, request); myErr != nil {
			myErr.HandleErr(writer)
		}
	}
}
