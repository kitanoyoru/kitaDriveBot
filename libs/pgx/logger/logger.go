package logger

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type Logger struct {
	logger      zerolog.Logger
	fromContext bool
}

type Option func(logger *Logger)

func FromContext() Option {
	return func(logger *Logger) {
		logger.fromContext = true
	}
}

func NewLogger(logger zerolog.Logger, options ...Option) *Logger {
	l := Logger{
		logger: logger,
	}
	l.init(options)
	return &l
}

func (pl *Logger) init(options []Option) {
	for _, opt := range options {
		opt(pl)
	}
}

func (pl *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	var zLevel zerolog.Level
	switch level {
	case tracelog.LogLevelNone:
		zLevel = zerolog.NoLevel
	case tracelog.LogLevelError:
		zLevel = zerolog.ErrorLevel
	case tracelog.LogLevelWarn:
		zLevel = zerolog.WarnLevel
	case tracelog.LogLevelInfo:
		zLevel = zerolog.DebugLevel
	case tracelog.LogLevelDebug:
		zLevel = zerolog.DebugLevel
	default:
		zLevel = zerolog.DebugLevel
	}

	var zCtx zerolog.Context
	if pl.fromContext {
		logger := zerolog.Ctx(ctx)
		zCtx = logger.With()
	} else {
		zCtx = pl.logger.With()
	}

	pgxlog := zCtx.Logger()
	event := pgxlog.WithLevel(zLevel)
	if event.Enabled() {
		if pl.fromContext {
			event.Str("module", "pgx")
		}
		event.Fields(data).Msg(msg)
	}
}
