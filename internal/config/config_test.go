package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kitanoyoru/kitaDriveBot/internal/config"
)

func TestLoadFromEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "" +
		"TELEGRAM_BOT_TOKEN=test-token\n" +
		"OWNER_TELEGRAM_ID=12345\n" +
		"GOOGLE_CLIENT_ID=client-id\n" +
		"GOOGLE_CLIENT_SECRET=client-secret\n" +
		"DEFAULT_FOLDER_PATH=Imports\n"

	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := godotenv.Overload(envPath); err != nil {
		t.Fatalf("Overload() error = %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.TelegramBotToken != "test-token" {
		t.Fatalf("TelegramBotToken = %q, want test-token", cfg.TelegramBotToken)
	}
	if cfg.OwnerTelegramID != 12345 {
		t.Fatalf("OwnerTelegramID = %d, want 12345", cfg.OwnerTelegramID)
	}
	if cfg.DefaultFolderPath != "Imports" {
		t.Fatalf("DefaultFolderPath = %q, want Imports", cfg.DefaultFolderPath)
	}
}

func TestValidateForRun(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		TelegramBotToken:   "token",
		OwnerTelegramID:    1,
		GoogleClientID:     "id",
		GoogleClientSecret: "secret",
	}

	if err := cfg.ValidateForRun(); err != nil {
		t.Fatalf("ValidateForRun() error = %v", err)
	}
}

func TestValidateForRun_MissingToken(t *testing.T) {
	t.Parallel()

	err := (config.Config{}).ValidateForRun()
	if err == nil {
		t.Fatal("ValidateForRun() error = nil, want error")
	}
}
