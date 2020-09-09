package myerror

import (
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/server/interface/response"
)

// MyErr エラーコードを追加した独自の型
type MyErr struct {
	// TODO: pkg/errorsを使いたい
	errMsg  error
	errCode int32
}

// NewMyErr MyErrの定義
func NewMyErr(errMsg error, errCode int32) *MyErr {
	myErr := &MyErr{
		errMsg:  errMsg,
		errCode: errCode,
	}
	return myErr
}

// HandleErr エラー時の処理
func (myErr *MyErr) HandleErr(writer http.ResponseWriter) {
	if myErr.errCode == http.StatusBadRequest {
		log.Println(myErr.errMsg)
		response.BadRequest(writer, myErr.errMsg.Error())
		return
	} else if myErr.errCode == http.StatusInternalServerError {
		log.Println(myErr.errMsg)
		response.InternalServerError(writer, myErr.errMsg.Error())
		return
	} else {
		// TODO: エラーコードが400,500以外の場合の処理考える
		errMsg := fmt.Sprintf("!! errorCode mistake. errCode: %d !!", myErr.errCode)
		panic(errMsg)
	}
}

func (myErr *MyErr) Error() string {
	return myErr.errMsg.Error()
}
