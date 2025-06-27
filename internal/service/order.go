package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/AndreyKuskov2/gophermart/pkg/validator"
	"go.uber.org/zap"
)

type GophermartGetOrderStorager interface {
	GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Orders, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]models.Orders, error)
}

type GophermartCreateOrderStorager interface {
	CreateNewOrder(ctx context.Context, order *models.Orders) error
}

type GophermartOrderService struct {
	getStorage    GophermartGetOrderStorager
	createStorage GophermartCreateOrderStorager
	log           *logger.Logger
}

func NewGophermartOrderService(getStorage GophermartGetOrderStorager, createStorage GophermartCreateOrderStorager, log *logger.Logger) *GophermartOrderService {
	return &GophermartOrderService{
		getStorage:    getStorage,
		createStorage: createStorage,
		log:           log,
	}
}

func (gs *GophermartOrderService) CreateNewOrderService(ctx context.Context, orderNumber string, userID string) error {
	if !validator.LuhnAlgorith(orderNumber) {
		gs.log.Log.Info(ErrNumberIsNotCorrect.Error(), zap.String("order_number", orderNumber))
		return ErrNumberIsNotCorrect
	}

	order, err := gs.getStorage.GetOrderByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	currentUser, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	if order != nil {
		if order.UserID == currentUser {
			return ErrOrderAlreadyExists
		} else {
			return ErrOrderAlreadyExistsForAnotherUser
		}
	}

	newOrder := &models.Orders{
		Number: orderNumber,
		Status: "NEW",
		UserID: currentUser,
	}
	if err := gs.createStorage.CreateNewOrder(ctx, newOrder); err != nil {
		return err
	}

	return nil
}

func (gs *GophermartOrderService) GetOrdersService(ctx context.Context, userID string) ([]models.Orders, error) {
	return gs.getStorage.GetOrdersByUserID(ctx, userID)
}
