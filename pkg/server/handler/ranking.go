package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// Rank rankデータ
type Rank struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	RankNum  int32  `json:"rank"`
	Score    int32  `json:"score"`
}

type rankingListResponse struct {
	Ranks []Rank `json:"ranks"`
}

// HandleRankingList ランキング情報取得
func HandleRankingList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// クエリ取得
		query := request.URL.Query()
		if len(query["start"]) != 1 {
			log.Println("the length of query must be one.")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		startNum, err := strconv.Atoi(query["start"][0])
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if startNum <= 0 {
			log.Println(errors.New("query'start' must be positive"))
			response.BadRequest(writer, fmt.Sprintf("query'start' must be positive"))
			return
		}
		users, err := model.SelectUsersByHighScore(startNum)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: 順位範囲外の処理
		if len(users) == 0 {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. rank=%d", startNum))
			return
		}

		rankingList := make([]Rank, len(users), len(users))
		for i, user := range users {
			rankingList[i] = Rank{
				UserID:   user.ID,
				UserName: user.Name,
				RankNum:  int32(startNum + i),
				Score:    user.HighScore,
			}
		}

		response.Success(writer, &rankingListResponse{Ranks: rankingList})
	}

}
