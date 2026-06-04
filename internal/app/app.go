package app

import (
	"context"
	"fmt"

	"github.com/kitanoyoru/kitaDriveBot/internal/adapter/googledrive"
	"github.com/kitanoyoru/kitaDriveBot/internal/adapter/oauth"
	"github.com/kitanoyoru/kitaDriveBot/internal/adapter/telegram"
	"github.com/kitanoyoru/kitaDriveBot/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/internal/platform/log"
	"github.com/kitanoyoru/kitaDriveBot/internal/usecase"
)

func RunAuth(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateForAuth(); err != nil {
		return err
	}

	_ = log.New(cfg.LogLevel)

	authService := oauth.NewService(cfg)
	return authService.RunInteractiveAuth(ctx)
}

func RunBot(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateForRun(); err != nil {
		return err
	}

	logger := log.New(cfg.LogLevel)

	authService := oauth.NewService(cfg)
	if !authService.Store().Exists() {
		return fmt.Errorf("google token not found at %s; run `bot auth` first", cfg.GoogleTokenPath)
	}

	tokenSource, err := authService.TokenSource(ctx)
	if err != nil {
		return fmt.Errorf("load google token: %w", err)
	}

	driveRepo, err := googledrive.NewRepository(ctx, tokenSource, cfg.DriveRootFolderID)
	if err != nil {
		return err
	}

	fetcher := telegram.NewFileFetcher(cfg.TelegramBotToken)
	uploadService := usecase.NewUploadService(fetcher, driveRepo)

	tgBot, err := telegram.NewBot(
		cfg.TelegramBotToken,
		cfg.OwnerTelegramID,
		cfg.DefaultFolderPath,
		uploadService,
		logger,
	)
	if err != nil {
		return err
	}

	tgBot.Start(ctx)
	return nil
}
