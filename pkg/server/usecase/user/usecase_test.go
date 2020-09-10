package user

import (
	"testing"

	"github.com/google/uuid"

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

func TestUseCase_GetUserByUserID(t *testing.T) {
	// request
	request := exampleUser.ID
	// response
	expected := exampleUser

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
	mockUserRepository.EXPECT().SelectUserByUserID(request).Return(expected, nil)

	usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
	actual, myErr := usecase.GetUserByUserID(request)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}

func TestUseCase_GetUserByAuthToken(t *testing.T) {
	// request
	request := exampleUser.AuthToken
	// response
	expected := exampleUser

	// モックの設定
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := mock_user.NewMockUserRepo(ctrl)
	mockUserRepository.EXPECT().SelectUserByAuthToken(request).Return(expected, nil)

	usecase := NewUserUseCase(mockUserRepository, uuid.NewRandom)
	actual, myErr := usecase.GetUserByAuthToken(request)
	assert.Equal(t, expected, actual)
	assert.Empty(t, myErr)
}

func createMockUUID() (uuid.UUID, error) {
	var tmp [16]byte = [16]byte{}
	return tmp, nil
}

func TestUseCase_RegisterUserFromUserName(t *testing.T) {
	exampleUser := &model.UserL{
		ID:        "00000000-0000-0000-0000-000000000000",
		AuthToken: "00000000-0000-0000-0000-000000000000",
		Name:      "example_name",
		HighScore: 0,
		Coin:      0,
	}

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
}

func TestUseCase_UpdateUserName(t *testing.T) {
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
}
