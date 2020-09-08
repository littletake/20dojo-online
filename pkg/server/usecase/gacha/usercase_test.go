package gacha

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"20dojo-online/pkg/server/domain/repository/collection_item/mock_collection_item"
	"20dojo-online/pkg/server/domain/repository/gacha_probability/mock_gacha_probability"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/domain/repository/user_collection_item/mock_user_collection_item"
	"20dojo-online/pkg/test/testdata"
)

// 対応表作成のため順番に注意

func TestUseCase_CreateCItemSlice(t *testing.T) {
	returnCItemSlice := testdata.ReturnCItemSlice

	// request: nil
	// response: nil

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collection_item.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_user_collection_item.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gacha_probability.NewMockGachaProbRepository(ctrl)
	// DBからのレスポンスを固定
	mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(returnCItemSlice, nil)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
	)
	myErr := usecase.CreateCItemSlice()
	assert.Empty(t, myErr)

}
func TestUseCase_CreateItemRatioSlice(t *testing.T) {
	returnGachaProbSlice := testdata.ReturnGachaProbSlice

	// request: nil
	// response: nil

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collection_item.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_user_collection_item.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gacha_probability.NewMockGachaProbRepository(ctrl)
	// DBからのレスポンスを固定
	mockGachaProbRepository.EXPECT().SelectAllGachaProb().Return(returnGachaProbSlice, nil)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
	)
	myErr := usecase.CreateItemRatioSlice()
	assert.Empty(t, myErr)
}
func TestUseCase_GetItems(t *testing.T) {
	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collection_item.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_user_collection_item.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gacha_probability.NewMockGachaProbRepository(ctrl)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
	)
	t.Run("1 time", func(t *testing.T) {
		// request
		requestTimes := 1
		// response
		expected := testdata.ExampleGachaResult1.CollectionID

		itemSlice := usecase.GetItems(int32(requestTimes))
		// 個数の確認
		assert.Len(t, itemSlice, requestTimes)
		// 要素の確認
		actual := itemSlice[0]
		assert.Equal(t, expected, actual)

	})

	t.Run("10 times", func(t *testing.T) {
		// request
		requestTimes := 10
		// response
		min, err := strconv.Atoi(testdata.ExampleGachaResult1.CollectionID)
		if err != nil {
			log.Fatal(err)
		}
		max, err := strconv.Atoi(testdata.ExampleGachaResult3.CollectionID)
		if err != nil {
			log.Fatal(err)
		}

		itemSlice := usecase.GetItems(int32(requestTimes))
		// 個数の確認
		assert.Len(t, itemSlice, requestTimes)
		// 要素の確認
		for _, item := range itemSlice {
			itemInt, err := strconv.Atoi(item)
			if err != nil {
				log.Fatal(err)
			}
			if !(itemInt >= min && itemInt <= max) {
				err := fmt.Errorf("invaild itemID. errorItem: %d", itemInt)
				log.Fatal(err)
			}
		}
	})
}

// }

// func TestUseCase_Gacha(t *testing.T) {
// 	// request
// 	// TODO: 回数は1と10の二通り試す
// 	requestTimes := int32(1)
// 	requestID := testdata.ExampleUser.ID
// 	// response
// 	expected := []*model.GachaResult{
// 		testdata.ExampleGachaResult1,
// 	}

// 	exampleUser := &model.UserL{
// 		ID:        testdata.ExampleUser.ID,
// 		AuthToken: testdata.ExampleUser.AuthToken,
// 		Name:      testdata.ExampleUser.Name,
// 		HighScore: testdata.ExampleUser.HighScore,
// 		Coin:      100,
// 	}
// 	exampleNewItemSlice := []*model.UserCollectionItem{}
// 	returnUser := &model.UserL{
// 		ID:        testdata.ExampleUser.ID,
// 		AuthToken: testdata.ExampleUser.AuthToken,
// 		Name:      testdata.ExampleUser.Name,
// 		HighScore: testdata.ExampleUser.HighScore,
// 		Coin:      100 - requestTimes,
// 	}
// 	returnUCItemSlice := testdata.ReturnUCItemSlice
// 	returnCItemSlice := testdata.ReturnCItemSlice
// 	returnGachaProbSlice := testdata.ReturnGachaProbSlice

// 	// モックの設定
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
// 	mockCItemRepository := mock_collection_item.NewMockCItemRepository(ctrl)
// 	mockUCItemRepository := mock_user_collection_item.NewMockUCItemRepository(ctrl)
// 	mockGachaProbRepository := mock_gacha_probability.NewMockGachaProbRepository(ctrl)
// 	// DBからのレスポンスを固定
// 	mockUserRepository.EXPECT().SelectUserByUserID(requestID).Return(exampleUser, nil)
// 	mockUCItemRepository.EXPECT().SelectUCItemSliceByUserID(exampleUser.ID).Return(returnUCItemSlice, nil)
// 	mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(returnCItemSlice, nil)
// 	mockGachaProbRepository.EXPECT().SelectAllGachaProb().Return(returnGachaProbSlice, nil)
// 	tx, _ := db.Conn.Begin()
// 	// assert.NoError(t, err)
// 	mockUCItemRepository.EXPECT().BulkInsertUCItemSlice(exampleNewItemSlice, tx).Return(nil)
// 	mockUserRepository.EXPECT().UpdateUserByUserInTx(returnUser, tx).Return(nil)

// 	usecase := NewGachaUseCase(
// 		mockUserRepository,
// 		mockCItemRepository,
// 		mockUCItemRepository,
// 		mockGachaProbRepository,
// 	)
// 	actual, myErr := usecase.Gacha(requestTimes, requestID)
// 	assert.Equal(t, expected, actual)
// 	assert.Empty(t, myErr)
// }
