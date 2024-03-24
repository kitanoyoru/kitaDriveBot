package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	PostgreSQLName = "postgres"
)

func ConnectToDB(ctx context.Context, cfg *DatabaseConfig) (*sqlx.DB, error) {
	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, err
	}

	return sqlx.ConnectContext(ctx, cfg.Name, dsn)
}

func getDSN(cfg *DatabaseConfig) (string, error) {
	switch cfg.Name {
	case PostgreSQLName:
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.Database,
		), nil
	default:
		return "", fmt.Errorf("Specified database doesn't support: %s", cfg.Name)
	}
}
