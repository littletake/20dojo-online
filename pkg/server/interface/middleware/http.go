package middleware

// import "net/http"

// // Get GETリクエストを処理する
// func (m *middleware) Get(apiFunc http.Handler) http.Handler {
// 	return httpMethod(apiFunc, http.MethodGet)
// }

// // Post POSTリクエストを処理する
// // func (m *middleware) Post(apiFunc http.HandlerFunc) http.HandlerFunc {
// // 	return httpMethod(apiFunc, http.MethodPost)
// // }

// // httpMethod 指定したHTTPメソッドでAPIの処理を実行する
// func httpMethod(apiFunc http.Handler, method string) http.Handler {
// 	return func(writer http.ResponseWriter, request *http.Request) {

// 		// CORS対応
// 		writer.Header().Add("Access-Control-Allow-Origin", "*")
// 		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

// 		// プリフライトリクエストは処理を通さない
// 		if request.Method == http.MethodOptions {
// 			return
// 		}
// 		// 指定のHTTPメソッドでない場合はエラー
// 		if request.Method != method {
// 			writer.WriteHeader(http.StatusMethodNotAllowed)
// 			writer.Write([]byte("Method Not Allowed"))
// 			return
// 		}

// 		// 共通のレスポンスヘッダを設定
// 		writer.Header().Add("Content-Type", "application/json")
// 		apiFunc(writer, request)
// 	}
// }
