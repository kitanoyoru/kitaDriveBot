package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken   string
	OwnerTelegramID    int64
	GoogleClientID     string
	GoogleClientSecret string
	GoogleTokenPath    string
	DriveRootFolderID  string
	DefaultFolderPath  string
	LogLevel           string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	ownerID, err := parseOwnerID(os.Getenv("OWNER_TELEGRAM_ID"))
	if err != nil {
		return Config{}, err
	}

	tokenPath := os.Getenv("GOOGLE_TOKEN_PATH")
	if tokenPath == "" {
		tokenPath = "./data/token.json"
	}

	defaultPath := os.Getenv("DEFAULT_FOLDER_PATH")
	if defaultPath == "" {
		defaultPath = "Telegram"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	cfg := Config{
		TelegramBotToken:   strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		OwnerTelegramID:    ownerID,
		GoogleClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
		GoogleClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		GoogleTokenPath:    tokenPath,
		DriveRootFolderID:  strings.TrimSpace(os.Getenv("DRIVE_ROOT_FOLDER_ID")),
		DefaultFolderPath:  defaultPath,
		LogLevel:           logLevel,
	}

	return cfg, nil
}

func (c Config) ValidateForRun() error {
	if c.TelegramBotToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if c.OwnerTelegramID == 0 {
		return errors.New("OWNER_TELEGRAM_ID is required")
	}
	if c.GoogleClientID == "" {
		return errors.New("GOOGLE_CLIENT_ID is required")
	}
	if c.GoogleClientSecret == "" {
		return errors.New("GOOGLE_CLIENT_SECRET is required")
	}
	return nil
}

func (c Config) ValidateForAuth() error {
	if c.GoogleClientID == "" {
		return errors.New("GOOGLE_CLIENT_ID is required")
	}
	if c.GoogleClientSecret == "" {
		return errors.New("GOOGLE_CLIENT_SECRET is required")
	}
	return nil
}

func parseOwnerID(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid OWNER_TELEGRAM_ID: %w", err)
	}

	return id, nil
}
