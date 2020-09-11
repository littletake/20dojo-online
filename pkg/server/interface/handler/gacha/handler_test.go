package gacha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	usecase "20dojo-online/pkg/server/usecase/gacha"
	"20dojo-online/pkg/server/usecase/gacha/mock_gacha"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

func Test_HandleGachaDraw(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}
	var exampleGachaResult1 = &usecase.GachaResult{
		CollectionID: "1001",
		ItemName:     "example1",
		Rarity:       int32(1),
		IsNew:        false,
	}
	var exampleGachaResult2 = &usecase.GachaResult{
		CollectionID: "1002",
		ItemName:     "example2",
		Rarity:       int32(2),
		IsNew:        true,
	}
	var exampleGachaResult3 = &usecase.GachaResult{
		CollectionID: "1003",
		ItemName:     "example3",
		Rarity:       int32(3),
		IsNew:        true,
	}
	var returnGachaResults = []*usecase.GachaResult{
		exampleGachaResult1,
	}
	var returnGachaResults10 = []*usecase.GachaResult{
		exampleGachaResult1,
		exampleGachaResult1,
		exampleGachaResult1,
		exampleGachaResult1,
		exampleGachaResult2,
		exampleGachaResult2,
		exampleGachaResult2,
		exampleGachaResult3,
		exampleGachaResult3,
		exampleGachaResult3,
	}
	t.Run("正常系(ガチャ1回)", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		requestTimes := int32(1)
		reqBody, err := json.Marshal(GachaDrawRequest{Times: requestTimes})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/gacha/draw", reqBodyBf)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		mockGachaUseCase.EXPECT().Gacha(requestTimes, exampleUser.ID).Return(returnGachaResults, nil)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gachaHandler.HandleGachaDraw()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusOK, "./testdata/handleGachaDrawRes.golden")
	})
	t.Run("正常系(ガチャ10回)", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		requestTimes := int32(10)
		reqBody, err := json.Marshal(GachaDrawRequest{Times: requestTimes})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/gacha/draw", reqBodyBf)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		mockGachaUseCase.EXPECT().Gacha(requestTimes, exampleUser.ID).Return(returnGachaResults10, nil)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gachaHandler.HandleGachaDraw()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusOK, "./testdata/handleGachaDrawRes10.golden")
	})
	t.Run("準正常系：requestBodyから情報取得失敗", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		req := httptest.NewRequest("POST", "/gacha/draw", nil)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gachaHandler.HandleGachaDraw()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponseHeader(t, res, http.StatusInternalServerError)
	})
	t.Run("準正常系：ガチャの回数が1，10以外", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		requestTimes := int32(2)
		reqBody, err := json.Marshal(GachaDrawRequest{Times: requestTimes})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/gacha/draw", reqBodyBf)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gachaHandler.HandleGachaDraw()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusBadRequest, "./testdata/errGachaTimes.golden")
	})
	t.Run("準正常系：コンテキストからuserID取得失敗", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		requestTimes := int32(1)
		reqBody, err := json.Marshal(GachaDrawRequest{Times: requestTimes})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/gacha/draw", reqBodyBf)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(gachaHandler.HandleGachaDraw())
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusInternalServerError, "./testdata/errGetUserIDFromContext.golden")
	})
	t.Run("準正常系(Gacha())：取得失敗", func(t *testing.T) {
		// リクエストの設定
		// TODO: ガチャの回数を1回と10回試す
		requestTimes := int32(1)
		reqBody, err := json.Marshal(GachaDrawRequest{Times: requestTimes})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/gacha/draw", reqBodyBf)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()
		expectErr := myerror.NewMyErr(
			fmt.Errorf("Internal Server Error"),
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGachaUseCase := mock_gacha.NewMockGachaUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		mockGachaUseCase.EXPECT().Gacha(requestTimes, exampleUser.ID).Return(returnGachaResults, expectErr)

		// テストの実行
		gachaHandler := NewGachaHandler(mockGachaUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gachaHandler.HandleGachaDraw()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusInternalServerError, "./testdata/errGacha.golden")
	})

}
