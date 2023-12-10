package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	httpAPI "github.com/kitanoyoru/kitaDriveBot/apps/service/internal/api/http/v0"
	kafkaAPI "github.com/kitanoyoru/kitaDriveBot/apps/service/internal/api/kafka/v0"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/service"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

const (
	defaultAppShutdownTimeout = 5 * time.Second

	defaultHttpServerReadTimeout  = 10 * time.Second
	defaultHttpServerWriteTimeout = 10 * time.Second

	defaultHttpServerMaxHeaderBytes = 1 << 20
)

type App struct {
	config *config.Config

	logger *logger.Logger

	service *service.Service

	kafkaServer *message.Router
	httpServer  *http.Server
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

func (app *App) Run() error {
	if err := app.initAndStartKafkaServer(); err != nil {
		return err
	}

	if err := app.initAndRunHTTPServer(); err != nil {
		return err
	}

	return nil
}

func (app *App) initAndStartKafkaServer() error {
	router, err := kafkaAPI.NewKafkaAPI(app.config.Kafka, app.service).GetRouter()
	if err != nil {
		return err
	}
	app.kafkaServer = router

	go func() {
		app.logger.Zap.Info("Starting Kafka server...")
		if err := router.Run(context.Background()); err != nil {
			app.logger.Zap.Sugar().Fatalf("Failed to start Kafka server: %v", err)
		}
	}()

	return nil
}

func (app *App) initAndRunHTTPServer() error {
	router, err := httpAPI.NewHTTPApi(app.service).GetRouter()
	if err != nil {
		return err
	}

	addr := app.config.Http.GetAddr()

	app.httpServer = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    defaultHttpServerReadTimeout,
		WriteTimeout:   defaultHttpServerWriteTimeout,
		MaxHeaderBytes: defaultHttpServerMaxHeaderBytes,
	}

	go func() {
		app.logger.Zap.Info("Starting HTTP server...")
		if err := app.httpServer.ListenAndServe(); err != nil {
			app.logger.Zap.Sugar().Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	return nil
}
