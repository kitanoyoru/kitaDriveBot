package gateway

import (
	"context"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userApi "github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/app/gateway/api/v1/user"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/config"
	userService "github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user/service"
	userStoragePostgres "github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user/storage/postgres"
	"github.com/kitanoyoru/kitaDriveBot/libs/app"
	"github.com/kitanoyoru/kitaDriveBot/libs/grpc/interceptor"
	"github.com/kitanoyoru/kitaDriveBot/libs/hasher/bcrypt"
	pgxZerolog "github.com/kitanoyoru/kitaDriveBot/libs/pgx/logger"
	pgxpoolTxLib "github.com/kitanoyoru/kitaDriveBot/libs/transactor/pgxpool"
	pb "github.com/kitanoyoru/kitaDriveBot/protos/gen/go/user/v1"
)

func NewApp(config config.Config) (app.App, error) {
	lLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse log level")
		lLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lLevel)

	dbPool, err := getDBPool(config.DatabaseConfig.ConnectionString, config.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get db pool")
	}

	transactor := pgxpoolTxLib.NewTransactor(dbPool)

	hasher := bcrypt.NewPasswordHasher(5)

	// user
	userStorage := userStoragePostgres.New(dbPool)
	userService := userService.New(userStorage, hasher, transactor)

	var (
		unaryLogOnEvents  []logging.LoggableEvent
		streamLogOnEvents []logging.LoggableEvent
	)
	if config.LogLevel == "debug" {
		unaryLogOnEvents = []logging.LoggableEvent{logging.StartCall, logging.PayloadReceived, logging.FinishCall}
		streamLogOnEvents = []logging.LoggableEvent{logging.StartCall, logging.FinishCall}
	} else {
		unaryLogOnEvents = []logging.LoggableEvent{logging.FinishCall}
		streamLogOnEvents = []logging.LoggableEvent{logging.FinishCall}
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		log.Err(err).Any("panic", p).Bytes("stack", debug.Stack()).Msg("recovered from panic")
		return status.Errorf(codes.Internal, "%s", p)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptor.Logger(log.Logger), logging.WithLogOnEvents(unaryLogOnEvents...)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptor.Logger(log.Logger), logging.WithLogOnEvents(streamLogOnEvents...)),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userApi.NewUserServiceServer(userService))

	return &gateway{
		config:     config,
		grpcServer: grpcServer,
	}, nil
}

type gateway struct {
	config config.Config

	grpcServer *grpc.Server
}

func (a *gateway) Run() error {
	listener, err := net.Listen("tcp", a.config.GRPCEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to listen tcp")
	}

	log.Info().Msg("grpc server listening on " + a.config.GRPCEndpoint)
	return a.grpcServer.Serve(listener)
}

func (a *gateway) Close() error {
	a.grpcServer.GracefulStop()
	return nil
}

func getDBPool(sqlConnectionString, logLevel string) (*pgxpool.Pool, error) {
	pgxPoolConfig, err := pgxpool.ParseConfig(sqlConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse pgx config")
	}

	traceLogLevel, err := tracelog.LogLevelFromString(logLevel)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse log level")
		traceLogLevel = tracelog.LogLevelInfo
	}
	pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxZerolog.NewLogger(log.Logger, pgxZerolog.FromContext()),
		LogLevel: traceLogLevel,
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgxpool")
	}

	return dbPool, nil
}
