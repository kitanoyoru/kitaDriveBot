package config

import "github.com/kitanoyoru/kitaDriveBot/libs/database"

type Config struct {
	Logger                  `json:"logger"`
	GrpcConfig              `json:"http"`
	database.DatabaseConfig `json:"database"`
}

type Logger struct {
	LogLevel string `json:"log_level"`
}

type RedisConfig struct {
	Url string `json:"url"`
}

type GrpcConfig struct {
	GRPCEndpoint string `json:"endpoint"`
}
