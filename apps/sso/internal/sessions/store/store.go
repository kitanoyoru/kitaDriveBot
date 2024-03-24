package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/sessions/store/repositories"
	"github.com/kitanoyoru/kitaDriveBot/libs/database"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

type StoreSession struct {
	db   *sqlx.DB
	User *repositories.User
}

func NewStoreSession(logger *logger.Logger, config *database.DatabaseConfig) (*StoreSession, error) {
	db, err := database.ConnectToDB(config)
	if err != nil {
		return nil, err
	}

	// TODO: migrate

	personRepository := repositories.NewUser(logger, db)

	return &StoreSession{
		db,
		personRepository,
	}, nil
}

func (s *StoreSession) Close() error {
	return s.db.Close()
}
