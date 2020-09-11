package ranking

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"20dojo-online/pkg/server/usecase/ranking/mock_ranking"

	"github.com/golang/mock/gomock"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/interface/myerror"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

func Test_HandleRankingList(t *testing.T) {
	var exampleUser1 = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 100,
		Coin:      0,
	}
	var exampleUser2 = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 10,
		Coin:      0,
	}
	var exampleUser3 = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 1,
		Coin:      0,
	}
	var exampleUsers = []*model.UserL{
		exampleUser1,
		exampleUser2,
		exampleUser3,
	}
	t.Run("正常系", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := 1
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(requestStartNum))
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)
		mockRankingUseCase.EXPECT().GetUsersByHighScore(int32(requestStartNum)).Return(exampleUsers, nil)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusOK, "./testdata/handleRankingListRes.golden")
	})
	t.Run("準正常系：クエリ条件不十分(クエリが複数)", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := 1
		tmp := 1
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(requestStartNum))
		q.Add("errQuery", strconv.Itoa(tmp))
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusBadRequest, "./testdata/errNotOneQuery.golden")
	})
	t.Run("準正常系：クエリ条件不十分(クエリが整数以外)", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := "test"
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", requestStartNum)
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponseHeader(t, res, http.StatusInternalServerError)
	})

	t.Run("準正常系：クエリ条件不十分(クエリが0以下の数)", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := 0
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(requestStartNum))
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusBadRequest, "./testdata/errNotPositiveQuery.golden")
	})

	t.Run("準正常系：クエリ条件不十分(クエリが0以下の数)", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := 0
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(requestStartNum))
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponseHeader(t, res, http.StatusBadRequest)
	})

	t.Run("準正常系(GetUsersByHighScore())：ユーザ取得に失敗", func(t *testing.T) {
		// リクエストレスポンスの設定
		requestStartNum := 1
		req := httptest.NewRequest("GET", "/ranking/list", nil)
		req.Header.Set("x-token", exampleUser1.AuthToken)
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(requestStartNum))
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()
		expectErr := myerror.NewMyErr(
			fmt.Errorf("Internal Server Error"),
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
		mockRankingUseCase := mock_ranking.NewMockRankingUseCase(ctrl)
		mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser1.AuthToken).Return(exampleUser1, nil)
		mockRankingUseCase.EXPECT().GetUsersByHighScore(int32(requestStartNum)).Return(nil, expectErr)

		// テストの実行
		rankingHandler := NewRankingHandler(mockRankingUseCase)
		m := middleware.NewMyMiddleware(mockUserUseCase)
		handle := m.Get(m.Authenticate(rankingHandler.HandleRankingList()))
		handle.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		test.AssertResponse(t, res, http.StatusInternalServerError, "./testdata/errGetUsersByHighScore.golden")
	})

}
