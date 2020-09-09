package collection

import (
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/collection"
)

// CollectionHandler Handlerのインターフェース
type CollectionHandler interface {
	HandleCollectionList() middleware.MyHandlerFunc
}

// collectionHandler usecaseとhandlerをつなぐもの
type collectionHandler struct {
	collectionUseCase usecase.CollectionUseCase
}

// NewCollectionHandler Userデータに関するHandler
func NewCollectionHandler(cu usecase.CollectionUseCase) CollectionHandler {
	return &collectionHandler{
		collectionUseCase: cu,
	}
}

// HandleCollectionList ガチャ実行
func (ch *collectionHandler) HandleCollectionList() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
		// collectionListResponse レスポンス形式
		type collectionListResponse struct {
			Collections []*usecase.CollectionItemResult `json:"collections"`
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
		collectionList, myErr := ch.collectionUseCase.GetCollectionSlice(userID)
		if myErr != nil {
			return myErr
		}
		response.Success(writer, &collectionListResponse{
			Collections: collectionList,
		})
		return myErr

	}
}
