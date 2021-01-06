package setting

import (
	"net/http"

	"20dojo-online/pkg/constant"
	m "20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/interface/response"
)

// TODO: swaggerのdocに反映

// SettingHandler インターフェース　handlerのメソッド一覧
type SettingHandler interface {
	HandleSettingGet() m.MyHandlerFunc
}

type settingHandler struct {
}

// NewSettingHandler Handler
func NewSettingHandler() SettingHandler {
	return &settingHandler{}
}

// HandleSettingGet ゲーム設定情報取得処理
func (sh *settingHandler) HandleSettingGet() m.MyHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) *myerror.MyErr {
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
		// myErr := myerror.NewMyErr(errors.New("testtesttest"), int32(500))
		// return myErr
		return nil
	}
}
