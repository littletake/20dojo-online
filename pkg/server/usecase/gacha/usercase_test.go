package gacha

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
	"20dojo-online/pkg/server/domain/repository/collectionitem/mock_collectionitem"
	"20dojo-online/pkg/server/domain/repository/gachaprobability/mock_gachaprobability"
	"20dojo-online/pkg/server/domain/repository/transaction/mock_transaction"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/domain/repository/usercollectionitem/mock_usercollectionitem"
	"20dojo-online/pkg/server/interface/myerror"
)

// 対応表作成のため順番に注意

func TestUseCase_CreateCItemSlice(t *testing.T) {
	// フラグの都合上，先に準正常系をテストする
	t.Run("準正常系(SelectAllCollectionItem()): infra層でエラー発生", func(t *testing.T) {
		// request: nil
		// response
		err := errors.New("Internal Server Error")
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
		mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
		mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
		mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)
		// DBからのレスポンスを固定
		mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(nil, err)

		usecase := NewGachaUseCase(
			mockUserRepository,
			mockCItemRepository,
			mockUCItemRepository,
			mockGachaProbRepository,
			int64(1),
			mockTxRepo,
		)
		myErr := usecase.CreateCItemSlice()
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("正常系", func(t *testing.T) {
		returnCItemSlice := ExampleCItemSlice
		// request: nil
		// response: nil

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
		mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
		mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
		mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)
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
	})

}
func TestUseCase_CreateItemRatioSlice(t *testing.T) {
	// 先に準正常系からテストする
	t.Run("準正常系(SelectAllGachaProb()): infra層でエラー発生", func(t *testing.T) {
		// request: nil
		// response
		err := errors.New("Internal Server Error")
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
		mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
		mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
		mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)

		// DBからのレスポンスを固定
		mockGachaProbRepository.EXPECT().SelectAllGachaProb().Return(nil, err)

		usecase := NewGachaUseCase(
			mockUserRepository,
			mockCItemRepository,
			mockUCItemRepository,
			mockGachaProbRepository,
			int64(1),
			mockTxRepo,
		)
		myErr := usecase.CreateItemRatioSlice()
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("正常系", func(t *testing.T) {
		returnGachaProbSlice := ExampleGachaProbSlice

		// request: nil
		// response: nil

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
		mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
		mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
		mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)

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
	})

}
func TestUseCase_GetItems(t *testing.T) {
	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)

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
		expected := ExampleCItem1.ItemID // "1001"

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
		min, err := strconv.Atoi(ExampleCItem1.ItemID) // "1001"
		if err != nil {
			log.Fatal(err)
		}
		max, err := strconv.Atoi(ExampleCItem3.ItemID) // "1003"
		if err != nil {
			log.Fatal(err)
		}

		itemSlice := usecase.GetItems(int32(requestTimes))
		// 個数の確認
		assert.Len(t, itemSlice, requestTimes)
		// 要素の確認
		// 具体的なアイテムIDの検証ではなくIDの上下限の間にあることを確認する
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

func TestUseCase_CreateGachaResults(t *testing.T) {
	// request
	requestID := ExampleUser.ID
	requestSlice := []string{
		ExampleCItem1.ItemID, // "1001"
		ExampleCItem2.ItemID, // "1002"
		ExampleCItem3.ItemID, // "1003"
	}
	requestHasGotItemMap := map[string]bool{
		ExampleCItem1.ItemID: true,  // "1001": true
		ExampleCItem2.ItemID: false, // "1002": false
		ExampleCItem3.ItemID: false, // "1003": false
	}
	// response
	expectedGachaResultSlice := []*GachaResult{
		ExampleGachaResult1,
		ExampleGachaResult2,
		ExampleGachaResult3,
	}
	expectedNewItemSlice := []*ucm.UserCollectionItem{
		NewItem1,
		NewItem2,
	}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
		mockTxRepo,
	)
	actualGachaResultSlice, actualNewItemSlice := usecase.CreateGachaResults(requestSlice, requestHasGotItemMap, requestID)
	assert.Equal(t, expectedGachaResultSlice, actualGachaResultSlice)
	assert.Equal(t, expectedNewItemSlice, actualNewItemSlice)
}

// TODO: BulkInsertANdUpdate()のテストも
func TestUseCase_BulkInsertAndUpdate(t *testing.T) {
	// request
	requestUser := ExampleUser
	requestNewItemSlice := []*ucm.UserCollectionItem{
		NewItem1,
		NewItem2,
	}
	// TODO: tx	の対応
	tx := sql.Tx{}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
	mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)

	mockUCItemRepository.EXPECT().BulkInsert(requestNewItemSlice, &tx).Return(nil)
	mockUserRepository.EXPECT().UpdateUserByUserInTx(requestUser, &tx).Return(nil)

	usecase := NewGachaUseCase(
		mockUserRepository,
		mockCItemRepository,
		mockUCItemRepository,
		mockGachaProbRepository,
		int64(1),
		mockTxRepo,
	)
	err := usecase.BulkInsertAndUpdate(requestNewItemSlice, requestUser, &tx)
	assert.Empty(t, err)
}

// func TestUseCase_Gacha(t *testing.T) {
// 	returnUCItemSlice := ExampleUCItemSlice
// 	// returnUser := ExampleUser
// 	tx := sql.Tx{}
// 	exampleFunc := func(f func(*sql.Tx) error) error {
// 		err := f(&tx)
// 		assert.NoError(t, err)
// 		return nil
// 	}

// 	// request
// 	// TODO: 回数は1と10の二通り試す
// 	requestTimes := int32(1)
// 	requestID := ExampleUser.ID
// 	// response
// 	expected := []*GachaResult{
// 		ExampleGachaResult1,
// 		ExampleGachaResult2,
// 		ExampleGachaResult3,
// 	}

// 	// モックの設定
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
// 	mockCItemRepository := mock_collectionitem.NewMockCollectionItemRepo(ctrl)
// 	mockUCItemRepository := mock_usercollectionitem.NewMockUserCollectionItemRepo(ctrl)
// 	mockGachaProbRepository := mock_gachaprobability.NewMockGachaProbRepo(ctrl)
// 	mockTxRepo := mock_transaction.NewMockTxRepo(ctrl)
// 	// DBからのレスポンスを固定
// 	mockUserRepository.EXPECT().SelectUserByUserID(ExampleUser.ID).Return(ExampleUser, nil)
// 	mockUCItemRepository.EXPECT().SelectSliceByUserID(ExampleUser.ID).Return(returnUCItemSlice, nil)
// 	mockTxRepo.EXPECT().Transaction(exampleFunc).Return(nil)

// 	usecase := NewGachaUseCase(
// 		mockUserRepository,
// 		mockCItemRepository,
// 		mockUCItemRepository,
// 		mockGachaProbRepository,
// 		int64(1),
// 		mockTxRepo,
// 	)
// 	actual, myErr := usecase.Gacha(requestTimes, requestID)
// 	assert.Empty(t, myErr)
// 	assert.Equal(t, expected, actual)
// }
