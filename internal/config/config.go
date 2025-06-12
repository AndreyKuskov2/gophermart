package config

import (
	"fmt"
	"strings"

	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecretToken       string `env:"JWT_TOKEN"`
}

func NewConfig(log *logger.Logger) (*Config, error) {
	var cfg Config

	pflag.StringVarP(&cfg.RunAddress, "run-address", "a", "localhost:8000", "run address")
	// TODO: Добавить обработчик пустой строки в db uri
	pflag.StringVarP(&cfg.DatabaseURI, "database-uri", "d", "", "database uri")
	pflag.StringVarP(&cfg.AccrualSystemAddress, "accrual-system-address", "r", "localhost:8080", "accrual system address")
	pflag.StringVarP(&cfg.JWTSecretToken, "jwt-token", "j", "some-secret-token", "jwt token")

	pflag.Parse()

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("Unknown flag: %v", arg)
		}
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to get environment variable value")
	}

	// if cfg.DatabaseURI == "" {
	// 	return nil, fmt.Errorf("database-uri is required")
	// }

	return &cfg, nil
}
