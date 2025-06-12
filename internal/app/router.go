package app

import (
	"github.com/AndreyKuskov2/gophermart/internal/app/middlewares"
	"github.com/AndreyKuskov2/gophermart/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *App) GophermartRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middlewares.LoggerMiddleware(app.Log))
	router.Use(middleware.Recoverer)

	h := handlers.NewGophermartHandlers(nil, app.Cfg, app.Log)

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/login", h.LoginUserHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Post("/orders", h.CreateNewOrderHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/orders", h.GetOrdersHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/balance", h.GetBalanceHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Post("/balance/withdraw", h.WithdrawBalanceHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/withdrawals", h.WithdrawAlsHandler)
	})

	return router
}
