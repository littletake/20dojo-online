package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// GachaDrawRequest リクエスト形式
type GachaDrawRequest struct {
	Times int32 `json:"times"`
}

// Result レスポンス形式の中身
type Result struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

// GachaDrawResponse レスポンス形式
type GachaDrawResponse struct {
	Results []Result `json:"results"`
}

// HandleGachaDraw ガチャ実行
func HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody GachaDrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		gachaTimes := requestBody.Times
		if gachaTimes != 1 {
			if gachaTimes != 10 {
				log.Println(errors.New("query'times' must be 1 or 10"))
				response.BadRequest(writer, fmt.Sprintf("query'times' must be 1 or 10. times=%d", gachaTimes))
				return
			}
		}

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

		// 必要枚数分のコインがあるかどうかを判定
		if user.Coin < constant.GachaCoinConsumption {
			log.Println(errors.New("user doesn't have enough coins"))
			response.BadRequest(writer, fmt.Sprintf("user doesn't have enough coins. userID=%s, coin=%d", userID, user.Coin))
			return
		}

		// ガチャ実行
		gachaProbList, err := model.SelectAllGachaProb()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: エラーメッセージを訂正する
		if len(gachaProbList) == 0 {
			log.Println(errors.New("error"))
			response.BadRequest(writer, fmt.Sprintf("error"))
			return
		}

		// TODO: 一度だけ実行するようにする．redisで実装できればいい
		// TODO: 命名考える
		data := []int32{}
		count := int32(0)
		for _, item := range gachaProbList {
			count += item.Ratio
			data = append(data, count)
		}
		// 乱数生成
		rand.Seed(time.Now().UnixNano())
		random := rand.Int31n(data[len(data)-1])

		// 適している番号を見つける
		// TODO: 当たっているかどうかを判定する関数作成
		i := 0
		for {
			if data[i] > random {
				break
			}
			i++
		}
		// fmt.Print(random, i, gachaProbList[i].ItemID)
		collectionItem, err := model.SelectCollectionItemByItemID(gachaProbList[i].CollectionItemID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		fmt.Print(collectionItem)

		// 持っているかどうか判定
		userCollectionItem, err := model.SelectUserCollectionItemByItemIDAndUserID(collectionItem.ItemID, userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: 別関数を作成する
		// TODO: 新規アイテムの場合DBに保存
		// TODO: コインの消費を追加
		// TODO: トランザクションをはる
		results := make([]Result, 1)
		if userCollectionItem == nil {
			results[0] = Result{
				CollectionID: collectionItem.ItemID,
				ItemName:     collectionItem.ItemName,
				Rarity:       collectionItem.Rarity,
				IsNew:        false,
			}
		} else {
			results[0] = Result{
				CollectionID: collectionItem.ItemID,
				ItemName:     collectionItem.ItemName,
				Rarity:       collectionItem.Rarity,
				IsNew:        true,
			}
		}
		// fmt.Print(userCollectionItem)

		response.Success(writer, &GachaDrawResponse{Results: results})
	}
}

// // DrawGacha ガチャ実行
// func DrawGacha() {
// 	gachaProbList, err := model.SelectAllGachaProb()
// 	if err != nil {
// 		log.Println(err)
// 		response.InternalServerError(writer, "Internal Server Error")
// 		return
// 	}
// 	// TODO: 訂正エラーメッセージ
// 	if len(gachaProbList) == 0 {
// 		log.Println(errors.New("error"))
// 		response.BadRequest(writer, fmt.Sprintf("error"))
// 		return
// 	}

// }
