package app

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kitanoyoru/kitaDriveBo/apps/service/internal/config"
	"github.com/kitanoyoru/kitaDriveBo/apps/service/internal/service"
	"github.com/kitanoyoru/kitaDriveBo/apps/service/pkg/logger"
)

type App struct {
	config *config.Config

	logger *logger.Logger

	service *service.Service

	kafkaServer  *message.Router
	metricServer *http.Server
}

func NewApp(config *config.Config) (*App, error) {
	app := App{}

	app.config = config

	logger, err := logger.NewLogger(config.Logger)
	if err != nil {
		return nil, err
	}
	app.logger = logger

	service, err := service.NewService(config.Google, logger)
	if err != nil {
		return nil, err
	}
	app.service = service

	return &app, nil
}
