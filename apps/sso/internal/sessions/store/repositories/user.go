package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

type User struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewUser(logger *logger.Logger, db *sqlx.DB) *User {
	return &User{}
}
