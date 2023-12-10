package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ZapLogger = iota
)

type Logger struct {
	Zap *zap.Logger
}

func NewLogger(config *LoggerConfig) (*Logger, error) {
	switch config.Type {
	case ZapLogger:
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
			Zap: logger,
		}, nil

	}

	return nil, errors.New("requested logger not found")
}
