package game

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
	"20dojo-online/pkg/server/usecase/game/mock_game"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

func Test_handleGameFinish(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 100,
		Coin:      0,
	}

	t.Run("正常系", func(t *testing.T) {
		// リクエストの設定
		requestScore := int32(1000)
		returnCoin := requestScore // UpdateCoinAndHighScore()で返ってくる値

		reqBody, err := json.Marshal(GameFinishRequest{Score: requestScore})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/game/finish", reqBodyBf)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGameUseCase := mock_game.NewMockGameUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		mockGameUseCase.EXPECT().UpdateCoinAndHighScore(exampleUser.ID, requestScore).Return(returnCoin, nil)
		// テストの実行
		gameHandler := NewGameHandler(mockGameUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gameHandler.HandleGameFinish()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusOK, "./testdata/handleGameFinishRes.golden")
	})
	t.Run("準正常系：コンテキストからuserID取得失敗", func(t *testing.T) {
		// リクエストの設定
		requestScore := int32(1000)

		reqBody, err := json.Marshal(GameFinishRequest{Score: requestScore})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/game/finish", reqBodyBf)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGameUseCase := mock_game.NewMockGameUseCase(ctrl)
		// テストの実行
		gameHandler := NewGameHandler(mockGameUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(gameHandler.HandleGameFinish())
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusInternalServerError, "./testdata/errGetUserIDFromContext.golden")
	})
	t.Run("準正常系：requestBodyから情報取得失敗", func(t *testing.T) {
		// リクエストの設定
		req := httptest.NewRequest("POST", "/game/finish", nil)
		req.Header.Set("x-token", exampleUser.AuthToken)
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockGameUseCase := mock_game.NewMockGameUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		// テストの実行
		gameHandler := NewGameHandler(mockGameUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gameHandler.HandleGameFinish()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponseHeader(t, res, http.StatusInternalServerError)
	})
	t.Run("準正常系(UpdateCoinAndHighScore)：更新失敗", func(t *testing.T) {
		// リクエストの設定
		requestScore := int32(1000)
		returnCoin := int32(0) // エラー値

		reqBody, err := json.Marshal(GameFinishRequest{Score: requestScore})
		assert.NoError(t, err)
		reqBodyBf := bytes.NewBuffer(reqBody)
		req := httptest.NewRequest("POST", "/game/finish", reqBodyBf)
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
		mockGameUseCase := mock_game.NewMockGameUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
		mockGameUseCase.EXPECT().UpdateCoinAndHighScore(exampleUser.ID, requestScore).Return(returnCoin, expectErr)
		// テストの実行
		gameHandler := NewGameHandler(mockGameUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Post(m.Authenticate(gameHandler.HandleGameFinish()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusInternalServerError, "./testdata/errUpdateCoinAndHighScore.golden")
	})

}
