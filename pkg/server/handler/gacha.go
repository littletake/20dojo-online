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
		if gachaTimes != constant.MinGachaTimes {
			if gachaTimes != constant.MaxGachaTimes {
				log.Println(errors.New("query'times' must be 1 or 10"))
				response.BadRequest(writer, fmt.Sprintf("requestBody'times' must be 1 or 10. times=%d", gachaTimes))
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
		if user.Coin < constant.GachaCoinConsumption*gachaTimes {
			log.Println(errors.New("user doesn't have enough coins"))
			response.BadRequest(writer, fmt.Sprintf("user doesn't have enough coins. userID=%s, coin=%d", userID, user.Coin))
			return
		}

		// table:gacha_probabilityの全件取得
		gachaProbList, err := model.SelectAllGachaProb()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// table:collection_itemの全件取得
		collectionItemList, err := model.SelectAllCollectionItem()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// userIDに適合したcollectionItemを取得
		userCollectionItemList, err := model.SelectUserCollectionItemListByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// ユーザが所持しているアイテムを示すmap
		// TODO: キャッシュで実現したい
		// [注意] 一回のガチャで同じものがでた時に追加するので長さは指定しない
		hasCollectionItemMap := make(map[string]bool)
		for _, userCollectionItem := range userCollectionItemList {
			itemID := userCollectionItem.CollectionItemID
			hasCollectionItemMap[itemID] = true
		}

		// TODO: 一度だけ実行するようにする．redisで実装できればいい
		// TODO: "data"の名称考える
		data := []int32{}
		count := int32(0)
		for _, item := range gachaProbList {
			count += item.Ratio
			data = append(data, count)
		}
		// 乱数生成
		// TODO: 乱数発生のseedも考える
		rand.Seed(time.Now().UnixNano())
		// gettingItemSlice 当てたアイテムのIDのlist
		gettingItemSlice := make([]string, gachaTimes)
		// TODO: [fix]for文の中
		for i := int32(0); i < gachaTimes; i++ {
			randomNum := rand.Int31n(data[len(data)-1])
			index := detectNumber(randomNum, data)
			gettingItemSlice[i] = collectionItemList[index].ItemID
		}
		// gettingItemMap := make(map[string]string, gachaTimes)
		// for i := int32(0); i < gachaTimes; i++ {
		// 	randomNum := rand.Int31n(data[len(data)-1])
		// 	index := detectNumber(randomNum, data)
		// 	gettingItemMap[i] = collectionItemList[index].ItemID
		// }

		// アイテムの照合
		// TODO: アイテムの保存
		resultList := make([]Result, gachaTimes)
		var newItemList model.UserCollectionItemList

		// TODO: テスト作成
		for i := int32(0); i < gachaTimes; i++ {
			for _, item := range collectionItemList {
				// TODO: ここもmapの方が早いかも
				if item.ItemID == gettingItemSlice[i] {
					// 既出itemの確認
					// fmt.Print(hasGot(i, gettingItemSlice))
					if hasCollectionItemMap[item.ItemID] {
						result := Result{
							CollectionID: item.ItemID,
							ItemName:     item.ItemName,
							Rarity:       item.Rarity,
							IsNew:        false,
						}
						resultList[i] = result
					} else {
						result := Result{
							CollectionID: item.ItemID,
							ItemName:     item.ItemName,
							Rarity:       item.Rarity,
							IsNew:        true,
						}
						resultList[i] = result
						// 当ガチャで同一のアイテムがあるかの確認
						hasCollectionItemMap[item.ItemID] = true
						// 登録
						newItem := model.UserCollectionItem{
							UserID:           userID,
							CollectionItemID: item.ItemID,
						}
						newItemList = append(newItemList, &newItem)
					}
				}
			}
		}

		if len(newItemList) != 0 {
			if err := model.BulkInsertUserCollectionItem(newItemList); err != nil {
				log.Println(err)
				response.InternalServerError(writer, "Internal Server Error")
				return
			}
		}

		// コインを消費
		user.Coin = user.Coin - constant.GachaCoinConsumption*gachaTimes
		err = model.UpdateUserByPrimaryKey(user)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: トランザクションをはる

		response.Success(writer, &GachaDrawResponse{Results: resultList})
	}
}

// detectNumber 適している番号を見つける
func detectNumber(random int32, data []int32) int32 {
	// TODO: 当たっているかどうかを判定する関数を作成すること
	num := int32(0)
	for {
		if data[num] > random {
			break
		}
		num++
	}
	return num
}

// hasGot 一回のガチャで同じものがあるかを判定
func hasGot(index int32, gettingItemSlice []string) bool {
	flag := false
	for i, gettingItem := range gettingItemSlice {
		if int32(i) == index {
			continue
		}
		if gettingItem == gettingItemSlice[index] {
			flag = true
		}
	}
	return flag
}

// // TODO: 関数の命名
// func hasItem(itemID string, userCollectionItemList *model.UserCollectionItemList) bool {
// 	flag := false
// 	i := 0
// 	for {
// 		if userCollectionItemList[i].CollectionItemID == itemID {
// 			flag = true
// 			break
// 		}
// 		i++
// 	}
// 	return flag
// }

// // detactItem アイテムの照合
// // TODO: numberの型注意
// func detectItem(number string, collectionItem interface) *Result{
// 	// TODO: 引数の修正
// 	for _, item  := range collectionItem {
// 		if item.ItemID == number {
// 			result := Result(
// 				CollectionID: item.ItemID,
// 				ItemName: item.ItemName,
// 				Rarity:       item.Rarity,
// 				IsNew:        false,
// 			)}}
// 	return &result
// }

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
