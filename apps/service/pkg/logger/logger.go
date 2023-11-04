package logger

import (
	"github.com/kitanoyoru/kitaDriveBo/apps/service/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	backend *zap.Logger
}

func NewLogger(config *config.LoggerConfig) (*Logger, error) {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	logger, err := zap.NewProduction(options...)
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &Logger{
		backend: logger,
	}, nil
}
