package handler

import (
	"fmt"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	"20dojo-online/pkg/server/usecase"
)

// CollectionHandler Handlerのインターフェース
type CollectionHandler interface {
	HandleCollectionList() http.HandlerFunc
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
func (ch collectionHandler) HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// collectionListResponse レスポンス形式
		type collectionListResponse struct {
			Collections []*model.CollectionItemResult `json:"collections"`
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
		collectionList, myErr := ch.collectionUseCase.GetCollectionSlice(userID)
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		response.Success(writer, &collectionListResponse{
			Collections: collectionList,
		})

	}
}
