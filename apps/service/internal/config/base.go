package config

type Config struct {
	Google *GoogleDriveConfig `json:"google,omitempty"`
	Kafka  *KafkaConfig       `json:"kafka,omitempty"`
	Logger *LoggerConfig      `json:"logger,omitempty"`
}

type KafkaConfig struct {
	BrokerList []string `json:"broker_list"`
	Topic      string   `json:"topic"`
	MaxRetry   int      `json:"max_retry"`
}

type GoogleDriveConfig struct {
	ApiKey string `json:"api_key,omitempty"`
}

type LoggerConfig struct {
	Level string `json:"level"`
}
