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

type PrometheusConfig struct {
	Port string `json:"port"`
}

type HttpConfig struct {
	Port string `json:"port"`
}

type GoogleDriveConfig struct {
	CredentialsPath string `json:"credentials_path"`
	TokenPath       string `json:"token_path"`
}

type LoggerConfig struct {
	Level string `json:"level"`
}
