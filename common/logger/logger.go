package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New create new logger
func New(env string) *zap.Logger {
	if env == "development" {
		logger, _ := zap.NewDevelopment(
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.AddCaller(),
		)
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}

// WithMethod adds method to logger
func WithMethod(logger *zap.Logger, method string) *zap.Logger {
	return logger.With(zap.String("method", method))
}
