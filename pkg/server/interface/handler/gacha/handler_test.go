package gacha

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
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
		exampleGachaResult2,
		exampleGachaResult3,
	}

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

}
