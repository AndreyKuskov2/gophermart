package service

import (
	"context"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/AndreyKuskov2/gophermart/pkg/validator"
	"go.uber.org/zap"
)

type GophermartWithdrawStorager interface {
	CreateWithdrawal(ctx context.Context, withdrawal *models.WithdrawBalance) error
	GetWithdrawalByUserID(ctx context.Context, userID string) ([]models.WithdrawBalance, error)
}

type GophermartWithdrawService struct {
	storage GophermartWithdrawStorager
	balance GophermartUserBalanceStorager
	log     *logger.Logger
}

func NewGophermartWithdrawService(storage GophermartWithdrawStorager, balance GophermartUserBalanceStorager, log *logger.Logger) *GophermartWithdrawService {
	return &GophermartWithdrawService{
		storage: storage,
		balance: balance,
		log:     log,
	}
}

func (gs *GophermartWithdrawService) WithdrawBalanceService(ctx context.Context, userID string, withdrawBalance *models.WithdrawBalanceRequest) error {
	if !validator.LuhnAlgorith(withdrawBalance.Order) {
		gs.log.Log.Info(ErrNumberIsNotCorrect.Error(), zap.String("order_number", withdrawBalance.Order))
		return ErrNumberIsNotCorrect
	}

	balance, err := gs.balance.GetUserBalance(ctx, userID)
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
	if err := gs.storage.CreateWithdrawal(ctx, withdrawal); err != nil {
		return err
	}

	return nil
}

func (gs *GophermartWithdrawService) GetWithdrawalService(ctx context.Context, userID string) ([]models.WithdrawBalance, error) {
	return gs.storage.GetWithdrawalByUserID(ctx, userID)
}
