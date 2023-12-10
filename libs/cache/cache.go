package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultDialTimeout = 10 * time.Second

	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second

	defaultPoolSize    = 12
	defaultPoolTimeout = 30 * time.Second
)

func ConnectToRedis(cfg *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		DialTimeout:  defaultDialTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		PoolSize:     defaultPoolSize,
		PoolTimeout:  defaultPoolTimeout,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		return nil, fmt.Errorf("Failed to initialize Redis connection: %+v", err)
	}

	return client, nil
}
