package main

import (
	"log"

	"github.com/AndreyKuskov2/gophermart/internal/app"
	"github.com/AndreyKuskov2/gophermart/internal/config"
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

	app := app.NewApp(cfg, logger)

	app.Run()
}
