package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PostgreSQLName = "postgres"
)

func ConnectToDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	driver, err := getDialector(cfg)
	if err != nil {
		return nil, err
	}

	return gorm.Open(driver, &gorm.Config{
		FullSaveAssociations: true,
	})
}

func getDialector(cfg *DatabaseConfig) (gorm.Dialector, error) {
	switch cfg.Name {
	case PostgreSQLName:
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.Database,
		)
		return postgres.New(postgres.Config{
			DSN: dsn,
		}), nil
	default:
		return nil, fmt.Errorf("Specified database doesn't support: %s", cfg.Name)
	}
}
