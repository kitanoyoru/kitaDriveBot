package repos

import (
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type User struct {
	logger *logger.Logger

	client *redis.Client
}

func NewUser(logger *logger.Logger, db *gorm.DB) *User {
	return &User{}
}
