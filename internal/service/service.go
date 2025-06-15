package service

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/AndreyKuskov2/gophermart/pkg/utils"
)

type IGophermartStorage interface {
	CreateUser(user models.UserCreditials) (int, error)
	GetUserByLogin(user models.UserCreditials) (int, error)

	GetOrderByNumber(orderNumber string) (*models.Orders, error)
	CreateNewOrder(order *models.Orders) error
	GetOrdersByUserID(userID string) ([]models.Orders, error)
	GetUserBalance(userID string) (*models.Balance, error)
	CreateWithdrawal(withdrawal *models.WithdrawBalance) error
	GetWithdrawalByUserID(userID string) ([]models.WithdrawBalance, error)
}

type GophermartService struct {
	storage IGophermartStorage
	log     *logger.Logger
}

func NewGophermartService(storage IGophermartStorage, log *logger.Logger) *GophermartService {
	return &GophermartService{
		storage: storage,
		log:     log,
	}
}

func (gs *GophermartService) RegisterUserService(user models.UserCreditials) (int, error) {
	return gs.storage.CreateUser(user)
}

func (gs *GophermartService) GetUserService(user models.UserCreditials) (int, error) {
	return gs.storage.GetUserByLogin(user)
}

func (gs *GophermartService) CreateNewOrderService(orderNumber string, userID string) error {
	if !utils.LuhnAlgorith(orderNumber) {
		gs.log.Log.Info(ErrNumberIsNotCorrect.Error())
		return ErrNumberIsNotCorrect
	}

	order, err := gs.storage.GetOrderByNumber(orderNumber)
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
	if err := gs.storage.CreateNewOrder(newOrder); err != nil {
		return err
	}

	return nil
}

func (gs *GophermartService) GetOrdersService(userID string) ([]models.Orders, error) {
	return gs.storage.GetOrdersByUserID(userID)
}

func (gs *GophermartService) GetUserBalanceService(userID string) (*models.Balance, error) {
	return gs.storage.GetUserBalance(userID)
}

func (gs *GophermartService) WithdrawBalanceService(userID string, withdrawBalance *models.WithdrawBalanceRequest) error {
	if !utils.LuhnAlgorith(withdrawBalance.Order) {
		gs.log.Log.Info(ErrNumberIsNotCorrect.Error())
		return ErrNumberIsNotCorrect
	}

	balance, err := gs.storage.GetUserBalance(userID)
	if err != nil {
		return err
	}

	if balance.Current < float64(withdrawBalance.Sum) {
		return ErrInvalidWithdrawSum
	}

	withdrawal := &models.WithdrawBalance{
		UserID:      userID,
		OrderNumber: withdrawBalance.Order,
		Amount:      withdrawBalance.Sum,
	}
	if err := gs.storage.CreateWithdrawal(withdrawal); err != nil {
		return err
	}

	return nil
}

func (gs *GophermartService) GetWithdrawalService(userID string) ([]models.WithdrawBalance, error) {
	return gs.storage.GetWithdrawalByUserID(userID)
}
