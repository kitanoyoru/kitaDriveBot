package app

import (
	"net"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/pkg/logger"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/api/grpc"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/config"
)

type App struct {
	config *config.Config

	logger *logger.Logger

	kafkaServer *message.Router
	httpServer  *http.Server
}

func (app *App) Run() error {
	if err := app.initAndRunGRPCServer(); err != nil {
		return err
	}

	return nil
}

func (app *App) initAndRunGRPCServer() error {
	lis, err := net.Listen("tcp", config.GrpcConfig.Port)
	if err != nil {
		return err
	}

	s, err := grpc.NewGrpcServer(app.config.Grpc)
	if err != nil {
		return err
	}

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
