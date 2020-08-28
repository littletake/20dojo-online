package handler

import (
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
)

// TODO: swaggerのdocに反映

// SettingGetResponse レスポンス形式
type SettingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
	RankingListNumber    int32 `json:"rankingListNumber"`
	MaxGachaTimes        int32 `json:"maxGachaTimes"`
	MinGachaTimes        int32 `json:"minGachaTimes"`
}

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response.Success(writer, &SettingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
			RankingListNumber:    constant.RankingListNumber,
			MaxGachaTimes:        constant.MaxGachaTimes,
			MinGachaTimes:        constant.MinGachaTimes,
		})
	}
}
