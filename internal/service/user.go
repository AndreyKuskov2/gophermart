package service

import (
	"context"

	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
)

type GophermartUserStorager interface {
	CreateUser(ctx context.Context, user models.UserCreditials) (int, error)
	GetUserByLogin(ctx context.Context, user models.UserCreditials) (int, error)
}

type GophermartUserService struct {
	storage GophermartUserStorager
	log     *logger.Logger
}

func NewGophermartUserService(storage GophermartUserStorager, log *logger.Logger) *GophermartUserService {
	return &GophermartUserService{
		storage: storage,
		log:     log,
	}
}

func (gs *GophermartUserService) RegisterUserService(ctx context.Context, user models.UserCreditials) (int, error) {
	return gs.storage.CreateUser(ctx, user)
}

func (gs *GophermartUserService) GetUserService(ctx context.Context, user models.UserCreditials) (int, error) {
	return gs.storage.GetUserByLogin(ctx, user)
}
