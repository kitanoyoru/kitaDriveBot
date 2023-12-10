package config

import (
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

type Config struct {
	Grpc   *GrpcConfig    `json:"http,omitempty"`
	Kafka  *KafkaConfig   `json:"kafka,omitempty"`
	Logger *logger.Logger `json:"logger,omitempty"`
}

type KafkaConfig struct {
	BrokerList []string `json:"broker_list"`
	Topic      string   `json:"topic"`
	MaxRetry   int      `json:"max_retry"`
}

type GrpcConfig struct {
	Port string `json:"port"`
}
