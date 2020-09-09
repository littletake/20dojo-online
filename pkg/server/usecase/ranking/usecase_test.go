package ranking

import (
	"testing"

	"20dojo-online/pkg/constant"
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

func TestUseCase_GetUsersByHighScore(t *testing.T) {
	// request
	request := int32(1)
	// response
	expected := []*model.UserL{
		exampleUser,
	}

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockModel := mock_user.NewMockUserRepo(ctrl)
	userMockModel.EXPECT().SelectUsersByHighScore(constant.RankingListNumber, request).Return(expected, nil)

	// usecase
	usecase := NewRankingUseCase(userMockModel)
	actual, myErr := usecase.GetUsersByHighScore(request)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}
