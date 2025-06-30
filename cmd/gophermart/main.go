package main

import (
	"context"
	"log"

	"github.com/AndreyKuskov2/gophermart/internal/app"
	"github.com/AndreyKuskov2/gophermart/internal/client"
	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/storage"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("cannot create logger")
	}

	cfg, err := config.NewConfig(logger)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	storage, err := storage.NewPostgres(cfg.DatabaseURI)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Info("migrations succesfully applied")

	accrualClient := client.NewClient(cfg.AccrualSystemAddress)

	accrualProcessor := app.NewAccrualProcessor(storage, accrualClient, logger)

	go accrualProcessor.Run(context.Background(), cfg.UpdateInterval, cfg.WorkerCount)

	app := app.NewApp(cfg, logger, storage)

	app.Run()
}
