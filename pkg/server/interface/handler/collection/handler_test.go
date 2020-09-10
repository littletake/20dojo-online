package collection

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/interface/middleware"
	usecase "20dojo-online/pkg/server/usecase/collection"
	"20dojo-online/pkg/server/usecase/collection/mock_collection"
	"20dojo-online/pkg/server/usecase/user/mock_user"
	"20dojo-online/pkg/test"
)

func Test_HandleColletionList(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 100,
		Coin:      0,
	}
	var exampleCItemResult1 = &usecase.CollectionItemResult{
		CollectionID: "1001",
		ItemName:     "example1",
		Rarity:       int32(1),
		HasItem:      true,
	}
	var exampleCItemResult2 = &usecase.CollectionItemResult{
		CollectionID: "1002",
		ItemName:     "example2",
		Rarity:       int32(2),
		HasItem:      false,
	}
	var exampleCItemResult3 = &usecase.CollectionItemResult{
		CollectionID: "1003",
		ItemName:     "example3",
		Rarity:       int32(3),
		HasItem:      false,
	}
	var returnCollectionItems = []*usecase.CollectionItemResult{
		exampleCItemResult1,
		exampleCItemResult2,
		exampleCItemResult3,
	}

	// リクエストの設定
	req := httptest.NewRequest("GET", "/collection/list", nil)
	req.Header.Set("x-token", exampleUser.AuthToken)
	rec := httptest.NewRecorder()

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := mock_user.NewMockUserUseCase(ctrl)
	mockCollectionUseCase := mock_collection.NewMockCollectionUseCase(ctrl)
	mockUserUseCase.EXPECT().GetUserByAuthToken(exampleUser.AuthToken).Return(exampleUser, nil)
	mockCollectionUseCase.EXPECT().GetCollectionSlice(exampleUser.ID).Return(returnCollectionItems, nil)
	// テストの実行
	collectionHandler := NewCollectionHandler(mockCollectionUseCase)
	m := middleware.NewMyMiddleware(mockUserUseCase)
	handle := m.Get(m.Authenticate(collectionHandler.HandleCollectionList()))
	handle.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	test.AssertResponse(t, res, http.StatusOK, "./testdata/handleCollectionListRes.golden")

}
