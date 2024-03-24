package gateway

import (
	"context"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
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
	"github.com/kitanoyoru/kitaDriveBot/libs/database"
	"github.com/kitanoyoru/kitaDriveBot/libs/grpc/interceptor"
	"github.com/kitanoyoru/kitaDriveBot/libs/hasher/bcrypt"
	sqlxTxLib "github.com/kitanoyoru/kitaDriveBot/libs/transactor/sqlx"
	pb "github.com/kitanoyoru/kitaDriveBot/protos/gen/go/user/v1"
)

func NewApp(config config.Config) (app.App, error) {
	ctx := context.Background()

	lLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse log level")
		lLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lLevel)

	db, err := database.ConnectToDB(ctx, config.DatabaseConfig)
	if err != nil {
		return nil, err
	}

	transactor := sqlxTxLib.NewTransactor(db)

	hasher := bcrypt.NewPasswordHasher(5)

	// user
	userStorage := userStoragePostgres.New(db)
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
