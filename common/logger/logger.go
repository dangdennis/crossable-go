package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

// GetLogger returns a logger instance
func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}
