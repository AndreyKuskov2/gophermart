package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() failed with error: %v", err)
	}

	if logger == nil {
		t.Fatal("NewLogger() returned nil logger")
	}

	if logger.Log == nil {
		t.Fatal("Logger.Log is nil")
	}

	logger.Log.Info("test message")
}

func TestLoggerStruct(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Log.Debug("debug message")
	logger.Log.Info("info message")
	logger.Log.Warn("warn message")
	logger.Log.Error("error message")

	logger.Log.Info("message with fields",
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
		zap.Bool("key3", true),
	)
}

func TestLoggerConfiguration(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Log.Info("info message should be logged")

	logger.Log.Debug("debug message should not be logged")
}
