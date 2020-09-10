package ranking

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"20dojo-online/pkg/server/usecase/ranking/mock_ranking"

	"github.com/golang/mock/gomock"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
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

}
