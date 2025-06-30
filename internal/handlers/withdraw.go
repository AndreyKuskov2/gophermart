package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/app/middlewares"
	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/internal/service"
	"github.com/AndreyKuskov2/gophermart/pkg/jwt"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type GophermartWithdrawServicer interface {
	WithdrawBalanceService(ctx context.Context, userID string, withdrawBalance *models.WithdrawBalanceRequest) error
	GetWithdrawalService(ctx context.Context, userID string) ([]models.WithdrawBalance, error)
}

type GophermartWithdrawHandlers struct {
	service GophermartWithdrawServicer
	cfg     *config.Config
	log     *logger.Logger
}

func NewGophermartWithdrawHandlers(service GophermartWithdrawServicer, cfg *config.Config, log *logger.Logger) *GophermartWithdrawHandlers {
	return &GophermartWithdrawHandlers{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (gh *GophermartWithdrawHandlers) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	var withdrawBalance models.WithdrawBalanceRequest
	if err := render.Bind(r, &withdrawBalance); err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	if err := gh.service.WithdrawBalanceService(r.Context(), claims.Subject, &withdrawBalance); err != nil {
		gh.log.Log.Info("failed to withdraw balance", zap.Error(err))
		switch {
		case errors.Is(err, service.ErrInvalidWithdrawSum):
			w.WriteHeader(http.StatusPaymentRequired)
			return
		case errors.Is(err, service.ErrNumberIsNotCorrect):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (gh *GophermartWithdrawHandlers) WithdrawAlsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	withdrawAls, err := gh.service.GetWithdrawalService(r.Context(), claims.Subject)
	if err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNoContent)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, withdrawAls)
}
