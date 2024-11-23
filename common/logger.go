package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// SetupLogger sets up the logger
func SetupLogger(env string) {
	if env == "development" {
		logger, _ = zap.NewDevelopment(
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	} else {
		logger, _ = zap.NewProduction()
	}
}

// Logger returns the logger
func Logger() *zap.Logger {
	if logger == nil {
		SetupLogger("development")
	}
	return logger
}
