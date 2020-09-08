package setting

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"20dojo-online/pkg/server/infra/persistence"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/usecase/user"
	"20dojo-online/pkg/test"
)

func TestSettingHandler(t *testing.T) {
	userPersistence := persistence.NewUserPersistence()
	userUseCase := user.NewUserUseCase(userPersistence)
	m := middleware.NewMiddleware(userUseCase)
	settingHandler := NewSettingHandler()

	req := httptest.NewRequest("GET", "/setting/get", nil)
	rec := httptest.NewRecorder()

	handle := m.Get(settingHandler.HandleSettingGet())
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	test.AssertResponse(t, res, http.StatusOK, "./testdata/response.golden")
}
