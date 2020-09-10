package ranking

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"20dojo-online/pkg/constant"
	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/domain/repository/user/mock_user"
	"20dojo-online/pkg/server/interface/myerror"

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
	t.Run("正常系", func(t *testing.T) {
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
	})

	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		request := int32(1)
		err := errors.New("Internal Server Error")
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userMockModel := mock_user.NewMockUserRepo(ctrl)
		userMockModel.EXPECT().SelectUsersByHighScore(constant.RankingListNumber, request).Return(nil, err)

		// usecase
		usecase := NewRankingUseCase(userMockModel)
		actual, myErr := usecase.GetUsersByHighScore(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系: アイテム数以上の順位が開始位置に指定された場合の処理", func(t *testing.T) {
		// request
		request := int32(100000)
		err := fmt.Errorf("user not found. rank=%d", request)
		// response
		expected := []*model.UserL{}
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)
		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userMockModel := mock_user.NewMockUserRepo(ctrl)
		userMockModel.EXPECT().SelectUsersByHighScore(constant.RankingListNumber, request).Return(expected, nil)

		// usecase
		usecase := NewRankingUseCase(userMockModel)
		actual, myErr := usecase.GetUsersByHighScore(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

}
