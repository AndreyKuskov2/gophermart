package service

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGophermartUserStorager is a mock implementation of GophermartUserStorager
type MockGophermartUserStorager struct {
	mock.Mock
}

func (m *MockGophermartUserStorager) CreateUser(ctx context.Context, user models.UserCreditials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockGophermartUserStorager) GetUserByLogin(ctx context.Context, user models.UserCreditials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func TestNewGophermartUserService(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
	assert.Equal(t, log, service.log)
}

func TestGophermartUserService_RegisterUserService_Success(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "testuser",
		Password: "testpassword",
	}

	expectedUserID := 123
	mockStorage.On("CreateUser", ctx, user).Return(expectedUserID, nil)

	userID, err := service.RegisterUserService(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_RegisterUserService_Error(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "testuser",
		Password: "testpassword",
	}

	expectedError := errors.New("user already exists")
	mockStorage.On("CreateUser", ctx, user).Return(0, expectedError)

	userID, err := service.RegisterUserService(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_GetUserService_Success(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "testuser",
		Password: "testpassword",
	}

	expectedUserID := 123
	mockStorage.On("GetUserByLogin", ctx, user).Return(expectedUserID, nil)

	userID, err := service.GetUserService(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_GetUserService_Error(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "nonexistentuser",
		Password: "wrongpassword",
	}

	expectedError := errors.New("user not found")
	mockStorage.On("GetUserByLogin", ctx, user).Return(0, expectedError)

	userID, err := service.GetUserService(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_RegisterUserService_EmptyCredentials(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "",
		Password: "",
	}

	expectedError := errors.New("invalid credentials")
	mockStorage.On("CreateUser", ctx, user).Return(0, expectedError)

	userID, err := service.RegisterUserService(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_GetUserService_EmptyCredentials(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx := context.Background()
	user := models.UserCreditials{
		Login:    "",
		Password: "",
	}

	expectedError := errors.New("invalid credentials")
	mockStorage.On("GetUserByLogin", ctx, user).Return(0, expectedError)

	userID, err := service.GetUserService(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, userID)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserService_WithNilLogger(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}

	service := NewGophermartUserService(mockStorage, nil)

	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
	assert.Nil(t, service.log)
}

func TestGophermartUserService_ContextCancellation(t *testing.T) {
	mockStorage := &MockGophermartUserStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserService(mockStorage, log)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	user := models.UserCreditials{
		Login:    "testuser",
		Password: "testpassword",
	}

	expectedError := context.Canceled
	mockStorage.On("CreateUser", ctx, user).Return(0, expectedError)

	userID, err := service.RegisterUserService(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, userID)
	mockStorage.AssertExpectations(t)
}
