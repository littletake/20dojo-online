package game

import (
	"testing"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var exampleUser = &model.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 0,
	Coin:      0,
}

func TestUseCase_UpdateCoinAndHighScore(t *testing.T) {
	// request
	requestID := exampleUser.ID
	requestScore := int32(10)
	// response
	expected := requestScore

	returnUser := &model.UserL{
		ID:        exampleUser.ID,
		AuthToken: exampleUser.AuthToken,
		Name:      exampleUser.Name,
		HighScore: requestScore,
		Coin:      exampleUser.Coin + requestScore,
	}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockModel := mock_user.NewMockUserRepo(ctrl)
	userMockModel.EXPECT().SelectUserByUserID(requestID).Return(exampleUser, nil)
	userMockModel.EXPECT().UpdateUserByUser(returnUser).Return(nil)

	usecase := NewGameUseCase(userMockModel)
	actual, myErr := usecase.UpdateCoinAndHighScore(requestID, requestScore)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}
