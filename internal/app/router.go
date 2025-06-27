package app

import (
	"github.com/AndreyKuskov2/gophermart/internal/app/middlewares"
	"github.com/AndreyKuskov2/gophermart/internal/handlers"
	"github.com/AndreyKuskov2/gophermart/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *App) GophermartRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middlewares.LoggerMiddleware(app.Log))
	router.Use(middleware.Recoverer)

	userService := service.NewGophermartUserService(app.Storage, app.Log)
	userHandlers := handlers.NewGophermartUserHandlers(userService, app.Cfg, app.Log)

	orderService := service.NewGophermartOrderService(app.Storage, app.Storage, app.Log)
	orderHandlers := handlers.NewGophermartOrderHandlers(orderService, app.Cfg, app.Log)

	balanceService := service.NewGophermartUserBalanceService(app.Storage, app.Log)
	balanceHandlers := handlers.NewGophermartBalanceHandlers(balanceService, app.Cfg, app.Log)

	withdrawService := service.NewGophermartWithdrawService(app.Storage, app.Storage, app.Log)
	withdrawHandlers := handlers.NewGophermartWithdrawHandlers(withdrawService, app.Cfg, app.Log)

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", userHandlers.RegisterUserHandler)
		r.Post("/login", userHandlers.LoginUserHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Post("/orders", orderHandlers.CreateNewOrderHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/orders", orderHandlers.GetOrdersHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/balance", balanceHandlers.GetBalanceHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Post("/balance/withdraw", withdrawHandlers.WithdrawBalanceHandler)
		r.With(middlewares.JwtAuthValidator(app.Cfg, app.Log)).Get("/withdrawals", withdrawHandlers.WithdrawAlsHandler)
	})

	return router
}
