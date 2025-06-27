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

// MockGophermartWithdrawStorager is a mock implementation of GophermartWithdrawStorager
type MockGophermartWithdrawStorager struct {
	mock.Mock
}

func (m *MockGophermartWithdrawStorager) CreateWithdrawal(ctx context.Context, withdrawal *models.WithdrawBalance) error {
	args := m.Called(ctx, withdrawal)
	return args.Error(0)
}

func (m *MockGophermartWithdrawStorager) GetWithdrawalByUserID(ctx context.Context, userID string) ([]models.WithdrawBalance, error) {
	args := m.Called(ctx, userID)
	withdrawals, _ := args.Get(0).([]models.WithdrawBalance)
	return withdrawals, args.Error(1)
}

func TestNewGophermartWithdrawService(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	assert.NotNil(t, service)
	assert.Equal(t, mockWithdrawStorage, service.storage)
	assert.Equal(t, mockBalanceStorage, service.balance)
	assert.Equal(t, log, service.log)
}

func TestGophermartWithdrawService_WithdrawBalanceService_Success(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   50.0,
	}

	userBalance := &models.Balance{
		Current:   100.0,
		Withdrawn: 25.0,
	}

	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(userBalance, nil)
	mockWithdrawStorage.On("CreateWithdrawal", ctx, mock.AnythingOfType("*models.WithdrawBalance")).Return(nil)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.NoError(t, err)
	mockBalanceStorage.AssertExpectations(t)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithdrawBalanceService_InvalidOrderNumber(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "1234567890", // invalid Luhn
		Sum:   50.0,
	}

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.ErrorIs(t, err, ErrNumberIsNotCorrect)
}

func TestGophermartWithdrawService_WithdrawBalanceService_InsufficientBalance(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   150.0,
	}

	userBalance := &models.Balance{
		Current:   100.0,
		Withdrawn: 25.0,
	}

	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(userBalance, nil)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.ErrorIs(t, err, ErrInvalidWithdrawSum)
	mockBalanceStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithdrawBalanceService_ExactBalance(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   100.0,
	}

	userBalance := &models.Balance{
		Current:   100.0,
		Withdrawn: 25.0,
	}

	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(userBalance, nil)
	mockWithdrawStorage.On("CreateWithdrawal", ctx, mock.AnythingOfType("*models.WithdrawBalance")).Return(nil)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.NoError(t, err)
	mockBalanceStorage.AssertExpectations(t)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithdrawBalanceService_GetBalanceError(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   50.0,
	}

	expectedError := errors.New("user not found")
	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.Equal(t, expectedError, err)
	mockBalanceStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithdrawBalanceService_CreateWithdrawalError(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   50.0,
	}

	userBalance := &models.Balance{
		Current:   100.0,
		Withdrawn: 25.0,
	}

	expectedError := errors.New("database error")
	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(userBalance, nil)
	mockWithdrawStorage.On("CreateWithdrawal", ctx, mock.AnythingOfType("*models.WithdrawBalance")).Return(expectedError)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.Equal(t, expectedError, err)
	mockBalanceStorage.AssertExpectations(t)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithdrawBalanceService_ZeroAmount(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713", // valid Luhn
		Sum:   0.0,
	}

	userBalance := &models.Balance{
		Current:   100.0,
		Withdrawn: 25.0,
	}

	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(userBalance, nil)
	mockWithdrawStorage.On("CreateWithdrawal", ctx, mock.AnythingOfType("*models.WithdrawBalance")).Return(nil)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.NoError(t, err)
	mockBalanceStorage.AssertExpectations(t)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_GetWithdrawalService_Success(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	expectedWithdrawals := []models.WithdrawBalance{
		{
			WithdrawalID: 1,
			UserID:       userID,
			OrderNumber:  "79927398713",
			Amount:       50.0,
		},
		{
			WithdrawalID: 2,
			UserID:       userID,
			OrderNumber:  "4532015112830366",
			Amount:       25.0,
		},
	}

	mockWithdrawStorage.On("GetWithdrawalByUserID", ctx, userID).Return(expectedWithdrawals, nil)

	withdrawals, err := service.GetWithdrawalService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWithdrawals, withdrawals)
	assert.Len(t, withdrawals, 2)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_GetWithdrawalService_EmptyList(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	expectedWithdrawals := []models.WithdrawBalance{}

	mockWithdrawStorage.On("GetWithdrawalByUserID", ctx, userID).Return(expectedWithdrawals, nil)

	withdrawals, err := service.GetWithdrawalService(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWithdrawals, withdrawals)
	assert.Len(t, withdrawals, 0)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_GetWithdrawalService_Error(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx := context.Background()
	userID := "123"
	expectedError := errors.New("database error")

	mockWithdrawStorage.On("GetWithdrawalByUserID", ctx, userID).Return(nil, expectedError)

	withdrawals, err := service.GetWithdrawalService(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, withdrawals)
	mockWithdrawStorage.AssertExpectations(t)
}

func TestGophermartWithdrawService_WithNilLogger(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, nil)

	assert.NotNil(t, service)
	assert.Equal(t, mockWithdrawStorage, service.storage)
	assert.Equal(t, mockBalanceStorage, service.balance)
	assert.Nil(t, service.log)
}

func TestGophermartWithdrawService_ContextCancellation(t *testing.T) {
	mockWithdrawStorage := &MockGophermartWithdrawStorager{}
	mockBalanceStorage := &MockGophermartUserBalanceStorager{}
	log, err := logger.NewLogger()
	assert.NoError(t, err)

	service := NewGophermartWithdrawService(mockWithdrawStorage, mockBalanceStorage, log)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	userID := "123"
	withdrawRequest := &models.WithdrawBalanceRequest{
		Order: "79927398713",
		Sum:   50.0,
	}

	expectedError := context.Canceled
	mockBalanceStorage.On("GetUserBalance", ctx, userID).Return(nil, expectedError)

	err = service.WithdrawBalanceService(ctx, userID, withdrawRequest)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockBalanceStorage.AssertExpectations(t)
}
