package repositories

import (
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
	"github.com/redis/go-redis/v9"
)

type Token struct {
	logger *logger.Logger

	client *redis.Client
}

func NewToken(logger *logger.Logger, config *config.RedisConfig) (*Token, error) {
	opts, err := redis.ParseURL(config.Url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	return &Token{
		logger,
		client,
	}, nil
}
