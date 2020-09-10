package setting

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

func TestSettingHandler(t *testing.T) {
	// リクエストレスポンスの設定
	req := httptest.NewRequest("GET", "/setting/get", nil)
	rec := httptest.NewRecorder()

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)

	// テストの実行
	m := middleware.NewMyMiddleware(mockUserUseCase)
	settingHandler := NewSettingHandler()
	handle := m.Get(settingHandler.HandleSettingGet())
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	test.AssertResponse(t, res, http.StatusOK, "./testdata/response.golden")
}
