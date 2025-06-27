package handlers

import (
	"context"
	"database/sql"
	"errors"
	"io"
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

type GophermartOrderServicer interface {
	CreateNewOrderService(ctx context.Context, orderNumber string, userID string) error
	GetOrdersService(ctx context.Context, userID string) ([]models.Orders, error)
}

type GophermartOrderHandlers struct {
	service GophermartOrderServicer
	cfg     *config.Config
	log     *logger.Logger
}

func NewGophermartOrderHandlers(service GophermartOrderServicer, cfg *config.Config, log *logger.Logger) *GophermartOrderHandlers {
	return &GophermartOrderHandlers{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (gh *GophermartOrderHandlers) CreateNewOrderHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}
	defer r.Body.Close()

	if err := gh.service.CreateNewOrderService(r.Context(), string(body), claims.Subject); err != nil {
		gh.log.Log.Info("failed to add order", zap.Error(err))
		switch {
		case errors.Is(err, service.ErrNumberIsNotCorrect):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, service.ErrOrderAlreadyExists):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, service.ErrOrderAlreadyExistsForAnotherUser):
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "")
}

func (gh *GophermartOrderHandlers) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	orders, err := gh.service.GetOrdersService(r.Context(), claims.Subject)
	if err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNoContent)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, orders)
}
