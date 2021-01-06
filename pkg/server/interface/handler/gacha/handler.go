package gacha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/gacha"
)

// GachaHandler UserにおけるHandlerのインターフェース
type GachaHandler interface {
	HandleGachaDraw() middleware.MyHandlerFunc
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

// GachaDrawRequest リクエスト形式
type GachaDrawRequest struct {
	Times int32 `json:"times"`
}

// GachaDrawResponse レスポンス形式
type GachaDrawResponse struct {
	Results []*usecase.GachaResult `json:"results"`
}

// HandleGachaDraw ガチャ実行
func (gh *gachaHandler) HandleGachaDraw() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
		// リクエストBodyから更新情報を取得
		var requestBody GachaDrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}
		// gachaTimes ガチャの回数
		gachaTimes := requestBody.Times
		if gachaTimes != constant.MinGachaTimes && gachaTimes != constant.MaxGachaTimes {
			myErr := myerror.NewMyErr(
				fmt.Errorf("requestBody'times' must be 1 or 10. times=%d", gachaTimes),
				http.StatusBadRequest,
			)
			return myErr
		}
		// コンテキストからuserID取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			myErr := myerror.NewMyErr(
				fmt.Errorf("userID is empty"),
				http.StatusInternalServerError,
			)
			return myErr
		}
		// 結果を取得
		gachaResultSlice, myErr := gh.gachaUseCase.Gacha(gachaTimes, userID)
		if myErr != nil {
			return myErr
		}
		response.Success(writer, &GachaDrawResponse{
			Results: gachaResultSlice,
		})
		return nil
	}
}
