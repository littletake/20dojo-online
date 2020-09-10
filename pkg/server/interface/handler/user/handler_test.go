package user

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"20dojo-online/pkg/server/infra/persistence"
// 	"20dojo-online/pkg/server/interface/middleware"
// 	"20dojo-online/pkg/server/usecase/user"
// 	"20dojo-online/pkg/test"
// )

// func TestHandler(t *testing.T) {
// 	userPersistence := persistence.NewUserPersistence()
// 	userUseCase := user.NewUserUseCase(userPersistence)
// 	userHandler := NewUserHandler(userUseCase)
// 	m := middleware.NewMiddleware(userUseCase)

// 	req := httptest.NewRequest("GET", "/user/get", nil)
// 	req.Header.Set("x-token", "ecfc137d-f211-4c72-9ea4-66bad81e62f0")
// 	rec := httptest.NewRecorder()

// 	handle := m.Get(m.Authenticate(userHandler.HandleUserGet()))
// 	handle.ServeHTTP(rec, req)
// 	res := rec.Result()
// 	defer res.Body.Close()
// 	fmt.Print(res)
// 	test.AssertResponse(t, res, http.StatusOK, "./testdata/response.golden")
// }
