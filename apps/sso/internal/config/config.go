package config

import (
	"github.com/kitanoyoru/kitaDriveBot/libs/cache"
	"github.com/kitanoyoru/kitaDriveBot/libs/database"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

type Config struct {
	Logger   *logger.Logger           `json:"logger,omitempty"`
	Grpc     *GrpcConfig              `json:"http,omitempty"`
	Database *database.DatabaseConfig `json:"database,omitempty"`
	Cache    *cache.RedisConfig       `json:"cache,omitempty"`
	Kafka    *KafkaConfig             `json:"kafka,omitempty"`
}

type KafkaConfig struct {
	BrokerList []string `json:"broker_list"`
	Topic      string   `json:"topic"`
	MaxRetry   int      `json:"max_retry"`
}

type RedisConfig struct {
	Url string `json:"url"`
}

type GrpcConfig struct {
	Port string `json:"port"`
}
