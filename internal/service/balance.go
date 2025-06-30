package service

import (
	"context"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
)

type GophermartUserBalanceStorager interface {
	GetUserBalance(ctx context.Context, userID string) (*models.Balance, error)
}

type GophermartUserBalanceService struct {
	storage GophermartUserBalanceStorager
	log     *logger.Logger
}

func NewGophermartUserBalanceService(storage GophermartUserBalanceStorager, log *logger.Logger) *GophermartUserBalanceService {
	return &GophermartUserBalanceService{
		storage: storage,
		log:     log,
	}
}

func (gs *GophermartUserBalanceService) GetUserBalanceService(ctx context.Context, userID string) (*models.Balance, error) {
	return gs.storage.GetUserBalance(ctx, userID)
}
