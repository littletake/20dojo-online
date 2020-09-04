package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// CollectionListResponse レスポンス形式
type CollectionListResponse struct {
	Collections []CollectionItem `json:"collections"`
}

// CollectionItem コレクションアイテム一覧
type CollectionItem struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

// HandleCollectionList コレクションアイテム一覧取得
func HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// ユーザデータの取得処理と存在チェックを実装
		user, err := model.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}
		// table:collection_itemの全件取得
		collectionItemSlice, err := model.SelectAllCollectionItem()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// table:user_collection_item からuserIDに適合したcollectionItemを取得
		userCollectionItemSlice, err := model.SelectUserCollectionItemSliceByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: キャッシュで実現したい
		// hasGotItemMap ユーザが所持しているアイテムを示すmap
		// [注意] ガチャ実行時も追加するので可変長指定
		hasGotItemMap := make(map[string]bool)
		for _, userCollectionItem := range userCollectionItemSlice {
			itemID := userCollectionItem.CollectionItemID
			hasGotItemMap[itemID] = true
		}

		// 二つのtableを合わせてresponseを作成
		collectionItemResult := make([]CollectionItem, len(collectionItemSlice))
		for i, collectionItem := range collectionItemSlice {
			result := CollectionItem{
				CollectionID: collectionItem.ItemID,
				ItemName:     collectionItem.ItemName,
				Rarity:       collectionItem.Rarity,
				HasItem:      hasGotItemMap[collectionItem.ItemID],
			}
			collectionItemResult[i] = result
		}

		response.Success(writer, &CollectionListResponse{
			Collections: collectionItemResult,
		})
	}
}
