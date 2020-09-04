package handler

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"20dojo-online/pkg/http/response"
// 	"20dojo-online/pkg/server/domain/model"
// )

// // DoErr エラー時の処理
// func DoErr(writer http.ResponseWriter, myErr *model.MyErr) {
// 	if myErr.ErrCode == 400 {
// 		log.Println(myErr.ErrMsg)
// 		response.BadRequest(writer, myErr.ErrMsg.Error())
// 		return
// 	} else if myErr.ErrCode == 500 {
// 		log.Println(myErr.ErrMsg)
// 		response.InternalServerError(writer, myErr.ErrMsg.Error())
// 		return
// 	} else {
// 		// TODO: エラーコードが400,500以外の場合の処理考える
// 		errMsg := fmt.Sprintf("!! errorCode mistake. ErrCode: %d !!", myErr.ErrCode)
// 		panic(errMsg)
// 	}

// }
