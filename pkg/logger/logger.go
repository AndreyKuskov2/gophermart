package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.Logger
}

func NewLogger() (*Logger, error) {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout", "logfile"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &Logger{
		Log: logger,
	}, nil
}
