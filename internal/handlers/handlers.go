package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/internal/storage"
	"github.com/AndreyKuskov2/gophermart/pkg/jwt"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
)

type IGophermartService interface {
	RegisterUserService(user models.UserCreditials) error
	GetUserService(user models.UserCreditials) error
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

	if err := gh.service.RegisterUserService(user); err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, storage.UserIsExist) {
			render.Status(r, http.StatusConflict)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.JWTSecretToken, user.Login)
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

	if err := gh.service.GetUserService(user); err != nil {
		gh.log.Log.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, storage.InvalidData) {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, "")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.JWTSecretToken, user.Login)
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
	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) WithdrawAlsHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}
