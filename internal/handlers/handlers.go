package handlers

import (
	"database/sql"
	"errors"
	"io"
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/app/middlewares"
	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/internal/service"
	"github.com/AndreyKuskov2/gophermart/internal/storage"
	"github.com/AndreyKuskov2/gophermart/pkg/jwt"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type IGophermartService interface {
	RegisterUserService(user models.UserCreditials) (int, error)
	GetUserService(user models.UserCreditials) (int, error)

	CreateNewOrderService(orderNumber string, userID string) error
	GetOrdersService(userID string) ([]models.Orders, error)
	GetUserBalanceService(userID string) (*models.Balance, error)
}

type GophermartHandlers struct {
	service IGophermartService
	cfg     *config.Config
	log     *logger.Logger
}

func NewGophermartHandlers(service IGophermartService, cfg *config.Config, log *logger.Logger) *GophermartHandlers {
	return &GophermartHandlers{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (gh *GophermartHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreditials

	if err := render.Bind(r, &user); err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	userID, err := gh.service.RegisterUserService(user)
	if err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, storage.ErrUserIsExist) {
			render.Status(r, http.StatusConflict)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.JWTSecretToken, userID)
	if err != nil {
		gh.log.Log.Info("cannot create jwt token")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	w.Header().Set("Authorization", jwtToken)
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "")
}

func (gh *GophermartHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreditials

	if err := render.Bind(r, &user); err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	userID, err := gh.service.GetUserService(user)
	if err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, storage.ErrInvalidData) {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.JWTSecretToken, userID)
	if err != nil {
		gh.log.Log.Info("cannot create jwt token")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	w.Header().Set("Authorization", jwtToken)
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "")
}

func (gh *GophermartHandlers) CreateNewOrderHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := gh.service.CreateNewOrderService(string(body), claims.Subject); err != nil {
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

func (gh *GophermartHandlers) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	orders, err := gh.service.GetOrdersService(claims.Subject)
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

func (gh *GophermartHandlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	if !ok {
		gh.log.Log.Info("cannot get jwt claims")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	balance, err := gh.service.GetUserBalanceService(claims.Subject)
	if err != nil {
		gh.log.Log.Info(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, balance)
}

func (gh *GophermartHandlers) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// claims, ok := r.Context().Value(middlewares.ContextClaims).(*jwt.JWTClaims)
	// if !ok {
	// 	gh.log.Log.Info("cannot get jwt claims")
	// 	render.Status(r, http.StatusBadRequest)
	// 	render.PlainText(w, r, "")
	// 	return
	// }

	// var withdrawBalance models.WithdrawBalanceRequest
	// if err := render.Bind(r, &withdrawBalance); err != nil {
	// 	gh.log.Log.Info("cannot parse body")
	// 	render.Status(r, http.StatusBadRequest)
	// 	render.PlainText(w, r, "")
	// 	return
	// }

	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) WithdrawAlsHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}
