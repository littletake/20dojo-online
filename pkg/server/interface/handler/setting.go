package handler

import (
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
)

// TODO: swaggerのdocに反映

// SettingHandler インターフェース　handlerのメソッド一覧
type SettingHandler interface {
	HandleSettingGet(writer http.ResponseWriter, request *http.Request)
}

type settingHandler struct {
}

// NewSettingHandler Handler
func NewSettingHandler() SettingHandler {
	return &settingHandler{}
}

// HandleSettingGet ゲーム設定情報取得処理
func (sh settingHandler) HandleSettingGet(writer http.ResponseWriter, request *http.Request) {
	// settingGetResponse レスポンス形式
	type settingGetResponse struct {
		GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
		RankingListNumber    int32 `json:"rankingListNumber"`
		MaxGachaTimes        int32 `json:"maxGachaTimes"`
		MinGachaTimes        int32 `json:"minGachaTimes"`
	}

	response.Success(writer, &settingGetResponse{
		GachaCoinConsumption: constant.GachaCoinConsumption,
		RankingListNumber:    constant.RankingListNumber,
		MaxGachaTimes:        constant.MaxGachaTimes,
		MinGachaTimes:        constant.MinGachaTimes,
	})
}
