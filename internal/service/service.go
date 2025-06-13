package service

import (
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
)

type IGophermartStorage interface {
	CreateUser(user models.UserCreditials) error
	GetUserByLogin(user models.UserCreditials) error
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

func (gs *GophermartService) RegisterUserService(user models.UserCreditials) error {
	return gs.storage.CreateUser(user)
}

func (gs *GophermartService) GetUserService(user models.UserCreditials) error {
	return gs.storage.GetUserByLogin(user)
}
