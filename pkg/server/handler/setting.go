package handler

import (
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
)

// SettingGetResponse ゲーム設定情報のレスポンス
type SettingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
	RankingListNumber    int32 `json:"RankingListNumber"`
}

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response.Success(writer, &SettingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
			RankingListNumber:    constant.RankingListNumber,
		})
	}
}
