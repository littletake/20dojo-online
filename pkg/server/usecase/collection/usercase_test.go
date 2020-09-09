package collection

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cm "20dojo-online/pkg/server/domain/model/collectionitem"
	um "20dojo-online/pkg/server/domain/model/user"
	ucm "20dojo-online/pkg/server/domain/model/usercollectionitem"
	"20dojo-online/pkg/server/domain/repository/collectionitem/mock_collectionitem"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/domain/repository/usercollectionitem/mock_usercollectionitem"
)

var exampleUser = &um.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 0,
	Coin:      0,
}

// ExampleCItemResult1 CollectionItemResult の例
var exampleCItemResult1 = &CollectionItemResult{
	CollectionID: exampleCItem1.ItemID,
	ItemName:     exampleCItem1.ItemName,
	Rarity:       exampleCItem1.Rarity,
	HasItem:      true,
}

// ExampleCItemResult2 CollectionItemResult の例
var exampleCItemResult2 = &CollectionItemResult{
	CollectionID: exampleCItem2.ItemID,
	ItemName:     exampleCItem2.ItemName,
	Rarity:       exampleCItem2.Rarity,
	HasItem:      false,
}

// ExampleCItemResult3 CollectionItemResult の例
var exampleCItemResult3 = &CollectionItemResult{
	CollectionID: exampleCItem3.ItemID,
	ItemName:     exampleCItem3.ItemName,
	Rarity:       exampleCItem3.Rarity,
	HasItem:      false,
}

// returnUCItemSlice user_collection_item の例
var returnUCItemSlice = []*ucm.UserCollectionItem{
	exampleUCItem1,
}
var exampleUCItem1 = &ucm.UserCollectionItem{
	UserID:           exampleUser.ID,
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

func TestUseCase_GetCollectionSlice(t *testing.T) {
	// request
	request := exampleUser.ID
	// response
	expected := []*CollectionItemResult{
		exampleCItemResult1,
		exampleCItemResult2,
		exampleCItemResult3,
	}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepository(ctrl)
	mockCItemRepository := mock_collectionitem.NewMockCItemRepository(ctrl)
	mockUCItemRepository := mock_usercollectionitem.NewMockUCItemRepository(ctrl)
	// DBからのレスポンスを固定
	mockUserRepository.EXPECT().SelectUserByUserID(request).Return(exampleUser, nil)
	mockUCItemRepository.EXPECT().SelectUCItemSliceByUserID(exampleUser.ID).Return(returnUCItemSlice, nil)
	mockCItemRepository.EXPECT().SelectAllCollectionItem().Return(returnCItemSlice, nil)

	usecase := NewCollectionUseCase(mockUserRepository, mockCItemRepository, mockUCItemRepository)
	actual, myErr := usecase.GetCollectionSlice(request)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}
