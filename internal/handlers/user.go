package handlers

import (
	"context"
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

type GophermartUserServicer interface {
	RegisterUserService(ctx context.Context, user models.UserCreditials) (int, error)
	GetUserService(ctx context.Context, user models.UserCreditials) (int, error)
}

type GophermartUserHandlers struct {
	service GophermartUserServicer
	cfg     *config.Config
	log     *logger.Logger
}

func NewGophermartUserHandlers(service GophermartUserServicer, cfg *config.Config, log *logger.Logger) *GophermartUserHandlers {
	return &GophermartUserHandlers{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (gh *GophermartUserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreditials

	if err := render.Bind(r, &user); err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	userID, err := gh.service.RegisterUserService(r.Context(), user)
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

func (gh *GophermartUserHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreditials

	if err := render.Bind(r, &user); err != nil {
		gh.log.Log.Info("cannot parse body")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}

	userID, err := gh.service.GetUserService(r.Context(), user)
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
