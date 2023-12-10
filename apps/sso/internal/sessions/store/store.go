package store

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/models"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/repos"
	"github.com/kitanoyoru/kitaDriveBot/libs/database"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
	"gorm.io/gorm"
)

type StoreSession struct {
	db   *gorm.DB
	User *repos.User
}

func NewStoreSession(logger *logger.Logger, config *database.DatabaseConfig) (*StoreSession, error) {
	db, err := database.ConnectToDB(config)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(models.User{})

	personRepository := repos.NewUser(logger, db)

	return &StoreSession{
		db,
		personRepository,
	}, nil
}

func (s *StoreSession) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
