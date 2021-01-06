package myerror

import (
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/server/interface/response"
)

// MyErr エラーコードを追加した独自の型
type MyErr struct {
	errBody error
	errCode int32
}

// NewMyErr MyErrの定義
func NewMyErr(errBody error, errCode int32) *MyErr {
	myErr := &MyErr{
		errBody: errBody,
		errCode: errCode,
	}
	return myErr
}

// HandleErr エラー時の処理
func (myErr *MyErr) HandleErr(writer http.ResponseWriter) {
	if myErr.errCode == http.StatusBadRequest {
		log.Println(myErr.errBody)
		response.BadRequest(writer, myErr.errBody.Error())
		return
	} else if myErr.errCode == http.StatusInternalServerError {
		log.Println(myErr.errBody)
		response.InternalServerError(writer, myErr.errBody.Error())
		return
	} else {
		// TODO: エラーコードが400,500以外の場合の処理考える
		errBody := fmt.Sprintf("!! invalid errorCode. errCode: %d !!", myErr.errCode)
		panic(errBody)
	}
}

func (myErr *MyErr) Error() string {
	return myErr.errBody.Error()
}

// GetErrCode errCodeの出力
func (myErr *MyErr) GetErrCode() int32 {
	return myErr.errCode
}
