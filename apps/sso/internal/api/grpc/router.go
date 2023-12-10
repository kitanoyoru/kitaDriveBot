package grpc

import (
	"context"
	"log/slog"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpcServer(logger *logger.Logger, config *config.GrpcConfig) (*grpc.Server, error) {
	loggingOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger.Zap)

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger.Zap, loggingOpts...),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	//ssov0.RegisterAuthServer(s, ssoHandler)

	return s, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zap.DebugLevel
	}

	return grpc_zap.DefaultCodeToLevel(code)
}
