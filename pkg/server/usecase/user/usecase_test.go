package user

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	model "20dojo-online/pkg/server/domain/model/user"
	"20dojo-online/pkg/server/domain/repository/mock/mock_user"
	"20dojo-online/pkg/server/interface/myerror"
)

func Test_GetUserByUserID(t *testing.T) {
	// データ
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	t.Run("正常系", func(t *testing.T) {
		// request
		request := exampleUser.ID
		// response
		expected := exampleUser

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(request).Return(expected, nil)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByUserID(request)
		assert.Equal(t, expected, actual)
		assert.Empty(t, myErr)
	})

	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		request := "aaa"
		err := errors.New("Internal Server Error")
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(request).Return(nil, err)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByUserID(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系：userが存在しない", func(t *testing.T) {
		// request
		request := "aaa"
		err := fmt.Errorf("user not found. userID=%s", request)
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(request).Return(nil, nil)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByUserID(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})
}

func Test_GetUserByAuthToken(t *testing.T) {
	// データ
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	t.Run("正常系", func(t *testing.T) {
		// request
		request := exampleUser.AuthToken
		// response
		expected := exampleUser

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByAuthToken(request).Return(expected, nil)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByAuthToken(request)
		assert.Equal(t, expected, actual)
		assert.Empty(t, myErr)
	})

	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		request := "aaa"
		err := errors.New("Internal Server Error")
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByAuthToken(request).Return(nil, err)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByAuthToken(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系: userが存在しない", func(t *testing.T) {
		// request
		request := "aaa"
		err := fmt.Errorf("user not found. token=%s", request)
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByAuthToken(request).Return(nil, nil)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.GetUserByAuthToken(request)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})
}

func createMockUUID() (uuid.UUID, error) {
	var tmp [16]byte = [16]byte{}
	return tmp, nil
}

func Test_RegisterUserFromUserName(t *testing.T) {
	exampleUser := &model.UserL{
		ID:        "00000000-0000-0000-0000-000000000000",
		AuthToken: "00000000-0000-0000-0000-000000000000",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	t.Run("正常系", func(t *testing.T) {
		// request
		requestName := exampleUser.Name
		// response
		expected := exampleUser.AuthToken

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().InsertUser(exampleUser).Return(nil)

		usecase := NewUserUseCase(mockUserRepository, createMockUUID)
		actual, myErr := usecase.RegisterUserFromUserName(requestName)
		assert.Equal(t, expected, actual)
		assert.Empty(t, myErr)
	})

	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		requestName := exampleUser.Name
		err := fmt.Errorf("Internal Server Error")
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().InsertUser(exampleUser).Return(err)

		usecase := NewUserUseCase(mockUserRepository, createMockUUID)
		actual, myErr := usecase.RegisterUserFromUserName(requestName)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})
}

func Test_UpdateUserName(t *testing.T) {
	var exampleUser = &model.UserL{
		ID:        "example_id",
		AuthToken: "example_token",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

	t.Run("正常系", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestName := "after"
		// response
		expected := &model.UserL{
			ID:        exampleUser.ID,
			AuthToken: exampleUser.AuthToken,
			Name:      requestName,
			HighScore: exampleUser.HighScore,
			Coin:      exampleUser.Coin,
		}

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(requestID).Return(exampleUser, nil)
		mockUserRepository.EXPECT().UpdateUserByUser(expected).Return(nil)

		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.UpdateUserName(requestID, requestName)
		assert.Equal(t, expected, actual)
		assert.Empty(t, myErr)
	})
	t.Run("準正常系: infra層でエラー発生", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestName := "after"
		err := errors.New("Internal Server Error")
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(requestID).Return(nil, err)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.UpdateUserName(requestID, requestName)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})

	t.Run("準正常系：userが存在しない", func(t *testing.T) {
		// request
		requestID := "aaa"
		requestName := "after"
		err := fmt.Errorf("user not found. userID=%s", requestID)
		// response
		expectErr := myerror.NewMyErr(
			err,
			http.StatusBadRequest,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(requestID).Return(nil, nil)
		// テストの実行
		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.UpdateUserName(requestID, requestName)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})
	t.Run("準正常系: 更新が失敗", func(t *testing.T) {
		// request
		requestID := exampleUser.ID
		requestName := "after"
		err := errors.New("Internal Server Error")
		// response
		expected := &model.UserL{
			ID:        exampleUser.ID,
			AuthToken: exampleUser.AuthToken,
			Name:      requestName,
			HighScore: exampleUser.HighScore,
			Coin:      exampleUser.Coin,
		}
		expectErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)

		// モックの設定
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepository := mock_user.NewMockUserRepo(ctrl)
		mockUserRepository.EXPECT().SelectUserByUserID(requestID).Return(exampleUser, nil)
		mockUserRepository.EXPECT().UpdateUserByUser(expected).Return(err)

		usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
		actual, myErr := usecase.UpdateUserName(requestID, requestName)
		assert.Empty(t, actual)
		assert.Equal(t, expectErr, myErr)
	})
}
