package handler

import (
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
)

// SettingGetResponse レスポンス形式
type SettingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
}

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response.Success(writer, &SettingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
		})
	}
}
