package handlers

import (
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
)

type IGophermartService interface {
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
	render.PlainText(w, r, "pong")
}

func (gh *GophermartHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
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
