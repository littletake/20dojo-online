package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"20dojo-online/pkg/server/model"

	"20dojo-online/pkg/http/response"
)

type rankingListResponse struct {
	Ranks interface{} `json:"ranks"`
}

// HandleRankingListGet ランキング情報取得
func HandleRankingListGet() http.HandlerFunc {
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
		rankList, err := model.SelectUsersByHighScore(startNum)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// TODO: [fix]順位範囲外の処理
		if len(rankList) == 0 {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. rank=%d", startNum))
			return
		}
		response.Success(writer, &rankingListResponse{Ranks: rankList})
	}

}
