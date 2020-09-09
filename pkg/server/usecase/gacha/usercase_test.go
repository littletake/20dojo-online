package gacha

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"20dojo-online/pkg/server/domain/repository/transaction/mock_transaction"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cm "20dojo-online/pkg/server/domain/model/collectionitem"
	gm "20dojo-online/pkg/server/domain/model/gachaprobability"
	um "20dojo-online/pkg/server/domain/model/user"
	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
	"20dojo-online/pkg/server/domain/repository/collectionitem/mock_collectionitem"
	"20dojo-online/pkg/server/domain/repository/gachaprobability/mock_gachaprobability"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/domain/repository/usercollectionitem/mock_usercollectionitem"
	crm "20dojo-online/pkg/server/usecase/collection"
)

// ExampleUser userL の例
var ExampleUser = &um.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 0,
	Coin:      0,
}

// ExampleCItemResult1 CollectionItemResult の例
var ExampleCItemResult1 = &crm.CollectionItemResult{
	CollectionID: exampleCItem1.ItemID,
	ItemName:     exampleCItem1.ItemName,
	Rarity:       exampleCItem1.Rarity,
	HasItem:      true,
}

// ExampleCItemResult2 CollectionItemResult の例
var ExampleCItemResult2 = &crm.CollectionItemResult{
	CollectionID: exampleCItem2.ItemID,
	ItemName:     exampleCItem2.ItemName,
	Rarity:       exampleCItem2.Rarity,
	HasItem:      false,
}

// ExampleCItemResult3 CollectionItemResult の例
var ExampleCItemResult3 = &crm.CollectionItemResult{
	CollectionID: exampleCItem3.ItemID,
	ItemName:     exampleCItem3.ItemName,
	Rarity:       exampleCItem3.Rarity,
	HasItem:      false,
}

// ReturnUCItemSlice user_collection_item の例
var ReturnUCItemSlice = []*ucm.UserCollectionItem{
	exampleUCItem1,
}
var exampleUCItem1 = &ucm.UserCollectionItem{
	UserID:           ExampleUser.ID,
	CollectionItemID: exampleCItem1.ItemID,
}

// returnCItemSlice collection_item の例
var returnCItemSlice = []*cm.CollectionItem{
	exampleCItem1,
	exampleCItem2,
	exampleCItem3,
}
var exampleCItem1 = &cm.CollectionItem{
	ItemID:   "1001",
	ItemName: "example1",
	Rarity:   int32(1),
}
var exampleCItem2 = &cm.CollectionItem{
	ItemID:   "1002",
	ItemName: "example2",
	Rarity:   int32(2),
}
var exampleCItem3 = &cm.CollectionItem{
	ItemID:   "1003",
	ItemName: "example3",
	Rarity:   int32(3),
}

// // ExampleGachaResultSlice gacha_probabiliy の例
// var ExampleGachaResultSlice = []*model.GachaResult{
// 	exampleGachaResult1,
// 	exampleGachaResult2,
// 	exampleGachaResult3,
// }

// exampleGachaResult1 GachaResult の例
var exampleGachaResult1 = &GachaResult{
	CollectionID: "1001",
	ItemName:     "example1",
	Rarity:       int32(1),
	IsNew:        false,
}

// ExampleGachaResult2 GachaResult の例
var ExampleGachaResult2 = &GachaResult{
	CollectionID: "1002",
	ItemName:     "example2",
	Rarity:       int32(2),
	IsNew:        false,
}

// exampleGachaResult3 GachaResult の例
var exampleGachaResult3 = &GachaResult{
	CollectionID: "1003",
	ItemName:     "example3",
	Rarity:       int32(3),
	IsNew:        false,
}

// returnGachaProbSlice GachaProb のスライスの例
var returnGachaProbSlice = []*gm.GachaProb{
	exampleGachaProb1,
	exampleGachaProb2,
	exampleGachaProb3,
}
var exampleGachaProb1 = &gm.GachaProb{
	CollectionItemID: "1001",
	Ratio:            6,
}
var exampleGachaProb2 = &gm.GachaProb{
	CollectionItemID: "1002",
	Ratio:            6,
}
var exampleGachaProb3 = &gm.GachaProb{
	CollectionItemID: "1003",
	Ratio:            6,
}

// 対応表作成のため順番に注意

func TestUseCase_CreateCItemSlice(t *testing.T) {
	// request: nil
	// response: nil

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepository(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepository(ctrl)
	// DBからのレスポンスを固定
	mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(returnCItemSlice, nil)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
		mockTxRepo,
	)
	myErr := usecase.CreateCItemSlice()
	assert.Empty(t, myErr)

}
func TestUseCase_CreateItemRatioSlice(t *testing.T) {
	// request: nil
	// response: nil

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepository(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepository(ctrl)

	// DBからのレスポンスを固定
	mockGachaProbRepository.EXPECT().SelectAllGachaProb().Return(returnGachaProbSlice, nil)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
		mockTxRepo,
	)
	myErr := usecase.CreateItemRatioSlice()
	assert.Empty(t, myErr)
}
func TestUseCase_GetItems(t *testing.T) {
	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUCItemRepository(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepository(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepository(ctrl)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
		mockTxRepo,
	)
	t.Run("1 time", func(t *testing.T) {
		// request
		requestTimes := 1
		// response
		expected := exampleGachaResult1.CollectionID

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
		min, err := strconv.Atoi(exampleGachaResult1.CollectionID)
		if err != nil {
			log.Fatal(err)
		}
		max, err := strconv.Atoi(exampleGachaResult3.CollectionID)
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
// 		testdata.exampleGachaResult1,
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
// 	returnCItemSlice := testdata.returnCItemSlice
// 	returnGachaProbSlice := testdata.returnGachaProbSlice

// 	// モックの設定
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
// 	mockCItemRepository := mock_collectionitem.NewMockCItemRepository(ctrl)
// 	mockUCItemRepository := mock_usercollectionitem.NewMockUCItemRepository(ctrl)
// 	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepository(ctrl)
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
