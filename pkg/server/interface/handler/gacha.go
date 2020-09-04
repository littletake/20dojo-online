package handler

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"time"

// 	"20dojo-online/pkg/constant"
// 	"20dojo-online/pkg/dcontext"
// 	"20dojo-online/pkg/http/response"
// 	"20dojo-online/pkg/server/domain/model"
// )

// // GachaHandler UserにおけるHandlerのインターフェース
// type GachaHandler interface {
// 	HandleGachaDraw(writer http.ResponseWriter, request *http.Request)
// }

// // gachaHandler usecaseとhandlerをつなぐもの
// type gachaHandler struct {
// 	gachaUseCase usecase.GachaUseCase
// }

// // NewGachaHandler Userデータに関するHandler
// func NewGachaHandler(gu usecase.GachaUseCase) GachaHandler {
// 	return &gachaHandler{
// 		gachaUseCase: gu,
// 	}
// }

// // init() 乱数のseed定義
// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// // HandleGachaDraw ガチャ実行
// func (gh gachaHandler) HandleGachaDraw(writer http.ResponseWriter, request *http.Request) {
// 	// gachaDrawRequest リクエスト形式
// 	type gachaDrawRequest struct {
// 		Times int32 `json:"times"`
// 	}

// 	// gachaDrawResponse レスポンス形式
// 	type gachaDrawResponse struct {
// 		Results []*model.GachaResult `json:"results"`
// 	}

// 	// リクエストBodyから更新情報を取得
// 	var requestBody gachaDrawRequest
// 	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
// 		myErr := usecase.CreateMyErr(
// 			err,
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// gachaTimes ガチャの回数
// 	gachaTimes := requestBody.Times
// 	if gachaTimes != constant.MinGachaTimes && gachaTimes != constant.MaxGachaTimes {
// 		myErr := usecase.CreateMyErr(
// 			fmt.Errorf("requestBody'times' must be 1 or 10. times=%d", gachaTimes),
// 			400,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// コンテキストからuserID取得
// 	ctx := request.Context()
// 	userID := dcontext.GetUserIDFromContext(ctx)
// 	if userID == "" {
// 		myErr := usecase.CreateMyErr(
// 			fmt.Errorf("userID is empty"),
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// 結果を取得
// 	gachaResultSlice, myErr := gh.gachaUseCase.Gacha(gachaTimes, userID)
// 	if myErr != nil {
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	response.Success(writer, &gachaDrawResponse{
// 		Results: gachaResultSlice,
// 	})

// }
