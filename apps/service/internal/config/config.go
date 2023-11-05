package config

import "fmt"

type Config struct {
	Google *GoogleDriveConfig `json:"google,omitempty"`
	Http   *HttpConfig        `json:"http,omitempty"`
	Kafka  *KafkaConfig       `json:"kafka,omitempty"`
	Logger *LoggerConfig      `json:"logger,omitempty"`
}

type KafkaConfig struct {
	BrokerList []string `json:"broker_list"`
	Topic      string   `json:"topic"`
	MaxRetry   int      `json:"max_retry"`
}

type HttpConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (hc *HttpConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", hc.Host, hc.Port)
}

type GoogleDriveConfig struct {
	CredentialsPath string `json:"credentials_path"`
	TokenPath       string `json:"token_path"`
}

type LoggerConfig struct {
	Type  int    `json:"type"`
	Level string `json:"level"`
}
