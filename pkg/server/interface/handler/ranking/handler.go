package ranking

import (
	"fmt"
	"net/http"
	"strconv"

	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/ranking"
)

// RankingHandler UserにおけるHandlerのインターフェース
type RankingHandler interface {
	HandleRankingList() http.HandlerFunc
}

// rankingHandler usecaseとhandlerをつなぐもの
type rankingHandler struct {
	rankingUseCase usecase.RankingUseCase
}

// NewRankingHandler Handlerを生成
func NewRankingHandler(ru usecase.RankingUseCase) RankingHandler {
	return &rankingHandler{
		rankingUseCase: ru,
	}
}

// HandleRankingList ランキング取得
func (rh *rankingHandler) HandleRankingList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// rankInfo ランキング情報
		type rankInfo struct {
			UserID   string `json:"userId"`
			UserName string `json:"userName"`
			Rank     int32  `json:"rank"`
			Score    int32  `json:"score"`
		}

		// RankingListResponse レスポンス形式
		type rankingListResponse struct {
			Ranks []rankInfo `json:"ranks"`
		}

		// クエリ取得
		query := request.URL.Query()
		if len(query["start"]) != 1 {
			myErr := myerror.NewMyErr(
				fmt.Errorf("the length of query must be one"),
				500,
			)
			myErr.HandleErr(writer)
			return
		}
		startNum, err := strconv.Atoi(query["start"][0])
		if err != nil {
			myErr := myerror.NewMyErr(err, 500)
			myErr.HandleErr(writer)
			return
		}
		if startNum <= 0 {
			myErr := myerror.NewMyErr(
				fmt.Errorf("query'start' must be positive"),
				400,
			)
			myErr.HandleErr(writer)
			return
		}
		// 対象範囲のユーザのスライス取得
		users, myErr := rh.rankingUseCase.GetUsersByHighScore(int32(startNum))
		if myErr != nil {
			myErr.HandleErr(writer)
			return
		}
		rankingList := make([]rankInfo, len(users), len(users))
		for i, user := range users {
			rankingList[i] = rankInfo{
				UserID:   user.ID,
				UserName: user.Name,
				Rank:     int32(startNum + i),
				Score:    user.HighScore,
			}
		}

		response.Success(writer, &rankingListResponse{
			Ranks: rankingList,
		})
	}
}
