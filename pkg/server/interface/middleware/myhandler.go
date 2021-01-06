package middleware

import (
	"net/http"

	"20dojo-online/pkg/server/interface/myerror"
)

// MyHandlerFunc 独自のハンドラ
// TODO: 独自の扱う型を定義する（MyErrの代わり）
// args1: statusCode
// args2: response（interface{}）
// args3: error <- 独自のエラー型にする？？
type MyHandlerFunc func(http.ResponseWriter, *http.Request) *myerror.MyErr

// // TODO: HandlerFuncの代わりになるようなHandlerをどのように作るか
// // ServeHTTP
// func (h MyHandlerFunc) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
// 	// API実行部分
// 	status, res, err := h(writer, request)
// 	fmt.Print(res)
// 	switch {
// 	// 200系
// 	case status == http.StatusOK:
// 		fmt.Println("OKOK")
// 		response.Success(writer, res)
// 		return
// 	// 400系
// 	// TODO: ログの保存方法も考える
// 	case status == http.StatusBadRequest:
// 		response.BadRequest(writer, err.Error())
// 		log.Print(res)
// 		return
// 		// 500系
// 		// TODO: ログの保存方法も考える
// 	case status == http.StatusInternalServerError:
// 		response.InternalServerError(writer, err.Error())
// 		log.Print(res)
// 		return
// 	}
// }
