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

// MockGophermartUserBalanceStorager is a mock implementation of GophermartUserBalanceStorager
type MockGophermartUserBalanceStorager struct {
	mock.Mock
}

func (m *MockGophermartUserBalanceStorager) GetUserBalance(ctx context.Context, userID string) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	balance, _ := args.Get(0).(*models.Balance)
	return balance, args.Error(1)
}

func TestNewGophermartUserBalanceService(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
	assert.Equal(t, log, service.log)
}

func TestGophermartUserBalanceService_GetUserBalanceService_Success(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "123"

	expectedBalance := &models.Balance{
		Current:   100.50,
		Withdrawn: 25.75,
	}

	mockStorage.On("GetUserBalance", ctx, userID).Return(expectedBalance, nil)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.Equal(t, 100.50, balance.Current)
	assert.Equal(t, float32(25.75), balance.Withdrawn)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_ZeroBalance(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "456"

	expectedBalance := &models.Balance{
		Current:   0.0,
		Withdrawn: 0.0,
	}

	mockStorage.On("GetUserBalance", ctx, userID).Return(expectedBalance, nil)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.Equal(t, 0.0, balance.Current)
	assert.Equal(t, float32(0.0), balance.Withdrawn)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_NegativeBalance(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "789"

	expectedBalance := &models.Balance{
		Current:   -50.25,
		Withdrawn: 100.0,
	}

	mockStorage.On("GetUserBalance", ctx, userID).Return(expectedBalance, nil)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.Equal(t, -50.25, balance.Current)
	assert.Equal(t, float32(100.0), balance.Withdrawn)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_Error(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "999"

	expectedError := errors.New("user not found")
	mockStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, balance)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_DatabaseError(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "111"

	expectedError := errors.New("database connection failed")
	mockStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, balance)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_EmptyUserID(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := ""

	expectedError := errors.New("invalid user ID")
	mockStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, balance)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_WithNilLogger(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}

	service := NewGophermartUserBalanceService(mockStorage, nil)

	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
	assert.Nil(t, service.log)
}

func TestGophermartUserBalanceService_ContextCancellation(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	userID := "123"

	expectedError := context.Canceled
	mockStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, balance)
	mockStorage.AssertExpectations(t)
}

func TestGophermartUserBalanceService_GetUserBalanceService_LargeNumbers(t *testing.T) {
	mockStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartUserBalanceService(mockStorage, log)

	ctx := context.Background()
	userID := "999999"

	expectedBalance := &models.Balance{
		Current:   999999.99,
		Withdrawn: 500000.50,
	}

	mockStorage.On("GetUserBalance", ctx, userID).Return(expectedBalance, nil)

	balance, err := service.GetUserBalanceService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.Equal(t, 999999.99, balance.Current)
	assert.Equal(t, float32(500000.50), balance.Withdrawn)
	mockStorage.AssertExpectations(t)
}
