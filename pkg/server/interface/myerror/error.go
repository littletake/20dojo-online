package myerror

import (
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/http/response"
)

// MyErr エラーコードを追加した独自の型
type MyErr struct {
	// TODO: pkg/errorsを使いたい
	ErrMsg  error
	ErrCode int32
}

// // CreateMyErr MyErrの定義
// func CreateMyErr(errMsg error, errCode int32) *MyErr {
// 	myErr := &MyErr{
// 		ErrMsg:  errMsg,
// 		ErrCode: errCode,
// 	}
// 	return myErr
// }

// HandleErr エラー時の処理
func (myErr MyErr) HandleErr(writer http.ResponseWriter) {
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
