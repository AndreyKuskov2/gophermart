package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/storage"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Cfg     *config.Config
	Log     *logger.Logger
	Storage *storage.Postgres
}

func NewApp(cfg *config.Config, log *logger.Logger, storage *storage.Postgres) *App {
	return &App{
		Cfg:     cfg,
		Log:     log,
		Storage: storage,
	}
}

func (app *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	router := app.GophermartRouter()

	go func() {
		app.Log.Log.Info("Start web-server", zap.String("address", app.Cfg.RunAddress))
		if err := http.ListenAndServe(app.Cfg.RunAddress, router); err != nil {
			app.Log.Log.Fatal("Failed to start server", zap.String("error", err.Error()))
		}
	}()

	<-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Log.Log.Info("Shutting down server...")
}
