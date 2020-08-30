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

// GachaDrawResponse レスポンス形式
type GachaDrawResponse struct {
	Results []GachaResult `json:"results"`
}

// GachaResult ガチャ結果
type GachaResult struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

// init() 乱数のseed定義
func init() {
	rand.Seed(time.Now().UnixNano())
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
		// gachaTimes ガチャの回数
		gachaTimes := requestBody.Times
		if gachaTimes != constant.MinGachaTimes {
			if gachaTimes != constant.MaxGachaTimes {
				log.Println(errors.New("requestBody'times' must be 1 or 10"))
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
		gachaProbSlice, err := model.SelectAllGachaProb()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// table:collection_itemの全件取得
		collectionItemSlice, err := model.SelectAllCollectionItem()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// userIDに適合したcollectionItemを取得
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

		// TODO: 初期化時に一度だけ実行したい．redisで実装??
		// itemRatioSlice ratioを考慮したアイテム対応表
		itemRatioSlice := make([]int32, len(gachaProbSlice))
		count := int32(0)
		for i, item := range gachaProbSlice {
			count += item.Ratio
			itemRatioSlice[i] = count
		}

		// 乱数によるガチャの実行
		// gettingItemSlice 当てたアイテムのIDのslice
		gettingItemSlice := make([]string, gachaTimes)
		for i := int32(0); i < gachaTimes; i++ {
			randomNum := rand.Int31n(itemRatioSlice[len(itemRatioSlice)-1])
			index := detectNumber(randomNum, itemRatioSlice)
			gettingItemSlice[i] = collectionItemSlice[index].ItemID
		}

		// アイテムの照合
		// TODO: アイテムの保存
		gachaResultSlice := make([]GachaResult, gachaTimes)
		var newItemSlice []*model.UserCollectionItem
		for i := int32(0); i < gachaTimes; i++ {
			for _, item := range collectionItemSlice {
				// TODO: 全探索改善する
				if item.ItemID == gettingItemSlice[i] {
					// 既出itemの確認
					// fmt.Print(hasGot(i, gettingItemSlice))
					if hasGotItemMap[item.ItemID] {
						result := GachaResult{
							CollectionID: item.ItemID,
							ItemName:     item.ItemName,
							Rarity:       item.Rarity,
							IsNew:        false,
						}
						gachaResultSlice[i] = result
					} else {
						result := GachaResult{
							CollectionID: item.ItemID,
							ItemName:     item.ItemName,
							Rarity:       item.Rarity,
							IsNew:        true,
						}
						gachaResultSlice[i] = result
						// 当ガチャで同一のアイテムがあるかの確認
						hasGotItemMap[item.ItemID] = true
						// 登録
						newItem := model.UserCollectionItem{
							UserID:           userID,
							CollectionItemID: item.ItemID,
						}
						newItemSlice = append(newItemSlice, &newItem)
					}
				}
			}
		}

		if len(newItemSlice) != 0 {
			if err := model.BulkInsertUserCollectionItem(newItemSlice); err != nil {
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
		response.Success(writer, &GachaDrawResponse{
			Results: gachaResultSlice,
		})
	}
}

// detectNumber 適している番号を見つける
func detectNumber(random int32, itemRatioSlice []int32) int32 {
	// TODO: 当たっているかどうかを判定する関数を作成すること
	num := int32(0)
	for {
		if itemRatioSlice[num] > random {
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
