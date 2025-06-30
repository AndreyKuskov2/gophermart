package handlers

import (
	"context"
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/app/middlewares"
	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/jwt"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
)

type GophermartBalanceServicer interface {
	GetUserBalanceService(ctx context.Context, userID string) (*models.Balance, error)
}

type GophermartBalanceHandlers struct {
	service GophermartBalanceServicer
	cfg     *config.Config
	log     *logger.Logger
}

func NewGophermartBalanceHandlers(service GophermartBalanceServicer, cfg *config.Config, log *logger.Logger) *GophermartBalanceHandlers {
	return &GophermartBalanceHandlers{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (gh *GophermartBalanceHandlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	balance, err := gh.service.GetUserBalanceService(r.Context(), claims.Subject)
	if err != nil {
		gh.log.Log.Info(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, balance)
}
