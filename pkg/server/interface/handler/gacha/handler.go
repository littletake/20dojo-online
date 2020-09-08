package gacha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/gacha"
)

// GachaHandler UserにおけるHandlerのインターフェース
type GachaHandler interface {
	HandleGachaDraw() http.HandlerFunc
}

// gachaHandler usecaseとhandlerをつなぐもの
type gachaHandler struct {
	gachaUseCase usecase.GachaUseCase
}

// NewGachaHandler Handlerを生成
func NewGachaHandler(gu usecase.GachaUseCase) GachaHandler {
	return &gachaHandler{
		gachaUseCase: gu,
	}
}

// HandleGachaDraw ガチャ実行
func (gh gachaHandler) HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// gachaDrawRequest リクエスト形式
		type gachaDrawRequest struct {
			Times int32 `json:"times"`
		}

		// gachaDrawResponse レスポンス形式
		type gachaDrawResponse struct {
			Results []*usecase.GachaResult `json:"results"`
		}

		// リクエストBodyから更新情報を取得
		var requestBody gachaDrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(err, 500)
			myErr.HandleErr(writer)
			return
		}
		// gachaTimes ガチャの回数
		gachaTimes := requestBody.Times
		if gachaTimes != constant.MinGachaTimes && gachaTimes != constant.MaxGachaTimes {
			myErr := myerror.NewMyErr(
				fmt.Errorf("requestBody'times' must be 1 or 10. times=%d", gachaTimes),
				400,
			)
			myErr.HandleErr(writer)
			return
		}
		// コンテキストからuserID取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("userID is empty"),
				500,
			)
			myErr.HandleErr(writer)
			return
		}
		// 結果を取得
		gachaResultSlice, myErr := gh.gachaUseCase.Gacha(gachaTimes, userID)
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		response.Success(writer, &gachaDrawResponse{
			Results: gachaResultSlice,
		})
	}
}
