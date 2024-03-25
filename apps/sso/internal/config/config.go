package config

import "github.com/kitanoyoru/kitaDriveBot/libs/database"

type Config struct {
	GrpcConfig              `json:"grpc"`
	database.DatabaseConfig `json:"database"`
	LoggerConfig            `json:"logger"`
}

type LoggerConfig struct {
	LogLevel string `json:"log_level"`
}

type GrpcConfig struct {
	GRPCEndpoint string `json:"endpoint"`
}
