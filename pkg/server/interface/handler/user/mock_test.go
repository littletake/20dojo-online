package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

// モックを使ったテスト
func Test_HandleUserGet(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}
	// リクエストレスポンスの設定
	req := httptest.NewRequest("GET", "/user/get", nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token", exampleUser.AuthToken)
	rec := httptest.NewRecorder()

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
	mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
	mockUserUseCase.EXPECT().GetUserByUserID(exampleUser.ID).Return(exampleUser, nil)

	// テストの実行
	userHandler := NewUserHandler(mockUserUseCase)
	m := middleware.NewMyMiddleware(mockUserUseCase)
	handle := m.Get(m.Authenticate(userHandler.HandleUserGet()))
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	test.AssertResponse(t, res, http.StatusOK, "./testdata/response.golden")
}

func Test_HandleUserCreate(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	// リクエストレスポンスの設定
	reqBody, err := json.Marshal(UserCreateRequest{Name: exampleUser.Name})
	assert.NoError(t, err)
	reqBodyBf := bytes.NewBuffer(reqBody)
	req := httptest.NewRequest("POST", "/user/create", reqBodyBf)
	req.Header.Set("x-token", exampleUser.AuthToken)
	rec := httptest.NewRecorder()

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
	mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
	mockUserUseCase.EXPECT().RegisterUserFromUserName(exampleUser.Name).Return(exampleUser.AuthToken, nil)

	// テストの実行
	userHandler := NewUserHandler(mockUserUseCase)
	m := middleware.NewMyMiddleware(mockUserUseCase)
	handle := m.Post(m.Authenticate(userHandler.HandleUserCreate()))
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	test.AssertResponse(t, res, http.StatusOK, "./testdata/handleUserCreateRes.golden")
}

func Test_HandleUserUpdate(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	// リクエストレスポンスの設定
	reqBody, err := json.Marshal(UserUpdateRequest{Name: exampleUser.Name})
	assert.NoError(t, err)
	reqBodyBf := bytes.NewBuffer(reqBody)
	req := httptest.NewRequest("POST", "/user/update", reqBodyBf)
	req.Header.Set("x-token", exampleUser.AuthToken)
	rec := httptest.NewRecorder()

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
	mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
	mockUserUseCase.EXPECT().UpdateUserName(exampleUser.ID, exampleUser.Name).Return(nil, nil)

	// テストの実行
	userHandler := NewUserHandler(mockUserUseCase)
	m := middleware.NewMyMiddleware(mockUserUseCase)
	handle := m.Post(m.Authenticate(userHandler.HandleUserUpdate()))
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	test.AssertResponseHeader(t, res, http.StatusOK)
}
