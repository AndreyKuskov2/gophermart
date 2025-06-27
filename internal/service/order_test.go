package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for GophermartGetOrderStorager
type MockGetOrderStorager struct {
	mock.Mock
}

func (m *MockGetOrderStorager) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Orders, error) {
	args := m.Called(ctx, orderNumber)
	order, _ := args.Get(0).(*models.Orders)
	return order, args.Error(1)
}

func (m *MockGetOrderStorager) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Orders, error) {
	args := m.Called(ctx, userID)
	orders, _ := args.Get(0).([]models.Orders)
	return orders, args.Error(1)
}

// Mock for GophermartCreateOrderStorager
type MockCreateOrderStorager struct {
	mock.Mock
}

func (m *MockCreateOrderStorager) CreateNewOrder(ctx context.Context, order *models.Orders) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestCreateNewOrderService_Success(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713" // valid Luhn
	userID := "1"

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(nil, sql.ErrNoRows)
	createStorage.On("CreateNewOrder", ctx, mock.AnythingOfType("*models.Orders")).Return(nil)

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.NoError(t, err)
	getStorage.AssertExpectations(t)
	createStorage.AssertExpectations(t)
}

func TestCreateNewOrderService_LuhnFail(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "1234567890" // invalid Luhn
	userID := "1"

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.ErrorIs(t, err, ErrNumberIsNotCorrect)
}

func TestCreateNewOrderService_OrderAlreadyExists_SameUser(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713"
	userID := "1"
	order := &models.Orders{Number: orderNumber, UserID: 1}

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(order, nil)

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.ErrorIs(t, err, ErrOrderAlreadyExists)
	getStorage.AssertExpectations(t)
}

func TestCreateNewOrderService_OrderAlreadyExists_AnotherUser(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713"
	userID := "1"
	order := &models.Orders{Number: orderNumber, UserID: 2}

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(order, nil)

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.ErrorIs(t, err, ErrOrderAlreadyExistsForAnotherUser)
	getStorage.AssertExpectations(t)
}

func TestCreateNewOrderService_GetOrderByNumberError(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713"
	userID := "1"
	someErr := errors.New("db error")

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(nil, someErr)

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.Equal(t, someErr, err)
	getStorage.AssertExpectations(t)
}

func TestCreateNewOrderService_CreateNewOrderError(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713"
	userID := "1"

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(nil, sql.ErrNoRows)
	createStorage.On("CreateNewOrder", ctx, mock.AnythingOfType("*models.Orders")).Return(errors.New("insert error"))

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.Error(t, err)
	createStorage.AssertExpectations(t)
}

func TestCreateNewOrderService_InvalidUserID(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	orderNumber := "79927398713"
	userID := "notanint"

	getStorage.On("GetOrderByNumber", ctx, orderNumber).Return(nil, sql.ErrNoRows)

	err := service.CreateNewOrderService(ctx, orderNumber, userID)
	assert.Error(t, err)
}

func TestGetOrdersService_Success(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	userID := "1"
	orders := []models.Orders{{OrderID: 1, Number: "79927398713", UserID: 1}}

	getStorage.On("GetOrdersByUserID", ctx, userID).Return(orders, nil)

	result, err := service.GetOrdersService(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, orders, result)
	getStorage.AssertExpectations(t)
}

func TestGetOrdersService_Error(t *testing.T) {
	getStorage := &MockGetOrderStorager{}
	createStorage := &MockCreateOrderStorager{}
	log, _ := logger.NewLogger()
	service := NewGophermartOrderService(getStorage, createStorage, log)

	ctx := context.Background()
	userID := "1"
	someErr := errors.New("db error")

	getStorage.On("GetOrdersByUserID", ctx, userID).Return(nil, someErr)

	result, err := service.GetOrdersService(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	getStorage.AssertExpectations(t)
}
