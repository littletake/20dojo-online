package collection

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository/collection_item/mock_collection_item"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/domain/repository/user_collection_item/mock_user_collection_item"
	"20dojo-online/pkg/test/testdata"
)

func TestUseCase_GetCollectionSlice(t *testing.T) {
	exampleUser := testdata.ExampleUser
	returnUCItemSlice := testdata.ReturnUCItemSlice
	returnCItemSlice := testdata.ReturnCItemSlice

	// request
	request := testdata.ExampleUser.ID
	// response
	expected := []*model.CollectionItemResult{
		testdata.ExampleCItemResult1,
		testdata.ExampleCItemResult2,
		testdata.ExampleCItemResult3,
	}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collection_item.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_user_collection_item.NewMockUCItemRepository(ctrl)
	// DBからのレスポンスを固定
	mockUserRepository.EXPECT().SelectUserByUserID(request).Return(exampleUser, nil)
	mockUCItemRepository.EXPECT().SelectUCItemSliceByUserID(exampleUser.ID).Return(returnUCItemSlice, nil)
	mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(returnCItemSlice, nil)

	usecase := NewCollectionUseCase(mockUserRepository, mockCItemRepository, mockUCItemRepository)
	actual, myErr := usecase.GetCollectionSlice(request)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}
