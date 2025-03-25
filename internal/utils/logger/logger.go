package logger

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	env := os.Getenv("APP_ENV") // Set APP_ENV to "production" or "development"
	var logger *zap.Logger
	var err error

	if env == "production" {
		logger, err = zap.NewProduction() // JSON format, info-level logs
	} else {
		logger, err = zap.NewDevelopment() // Human-readable format, debug-level logs
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}
