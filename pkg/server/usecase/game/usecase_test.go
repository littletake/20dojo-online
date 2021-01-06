package game

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/domain/repository/mock/mock_user"
	"20dojo-online/pkg/server/interface/myerror"
)

var exampleUser = &model.UserL{
	ID:        "example_id",
	AuthToken: "example_token",
	Name:      "example_name",
	HighScore: 0,
	Coin:      0,
}

func TestUseCase_UpdateCoinAndHighScore(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
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
	})

	t.Run("準正常系: scoreが自然数以外", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestScore := int32(-1)
		// response
		err := fmt.Errorf("score must be positive. score=%d", requestScore)
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userMockModel := mock_user.NewMockUserRepo(ctrl)

		usecase := NewGameUseCase(userMockModel)
		actual, myErr := usecase.UpdateCoinAndHighScore(requestID, requestScore)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)

	})

	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestScore := int32(10)
		// response
		err := errors.New("Internal Server Error")
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userMockModel := mock_user.NewMockUserRepo(ctrl)
		userMockModel.EXPECT().SelectUserByUserID(requestID).Return(nil, err)

		usecase := NewGameUseCase(userMockModel)
		actual, myErr := usecase.UpdateCoinAndHighScore(requestID, requestScore)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系：userが存在しない", func(t *testing.T) {
		// request
		requestID := "aaa"
		requestScore := int32(10)
		// response
		err := fmt.Errorf("user not found. userID=%s", requestID)
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userMockModel := mock_user.NewMockUserRepo(ctrl)
		userMockModel.EXPECT().SelectUserByUserID(requestID).Return(nil, nil)

		usecase := NewGameUseCase(userMockModel)
		actual, myErr := usecase.UpdateCoinAndHighScore(requestID, requestScore)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系: 更新に失敗", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestScore := int32(10)
		// response
		err := errors.New("Internal Server Error")
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

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
		userMockModel.EXPECT().UpdateUserByUser(returnUser).Return(err)

		usecase := NewGameUseCase(userMockModel)
		actual, myErr := usecase.UpdateCoinAndHighScore(requestID, requestScore)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

}
