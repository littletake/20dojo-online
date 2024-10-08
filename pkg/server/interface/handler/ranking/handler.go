package ranking

import (
	"fmt"
	"net/http"
	"strconv"

	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
	usecase "20dojo-online/pkg/server/usecase/ranking"
)

// RankingHandler UserにおけるHandlerのインターフェース
type RankingHandler interface {
	HandleRankingList() middleware.MyHandlerFunc
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
func (rh *rankingHandler) HandleRankingList() middleware.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
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
		if len(query) != 1 {
			myErr := myerror.NewMyErr(
				fmt.Errorf("invalid query. the length of query must be one"),
				http.StatusBadRequest,
			)
			return myErr
		}
		startNum, err := strconv.Atoi(query["start"][0])
		if err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return myErr
		}
		if startNum <= 0 {
			myErr := myerror.NewMyErr(
				fmt.Errorf("query'start' must be positive"),
				http.StatusBadRequest,
			)
			return myErr
		}
		// 対象範囲のユーザのスライス取得
		users, myErr := rh.rankingUseCase.GetUsersByHighScore(int32(startNum))
		if myErr != nil {
			return myErr
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
		return nil
	}
}
