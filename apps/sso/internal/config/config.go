package config

type Config struct {
	GrpcConfig     `json:"grpc"`
	DatabaseConfig `json:"database"`
	LoggerConfig   `json:"logger"`
}

type GrpcConfig struct {
	GRPCEndpoint string `json:"endpoint"`
}

type DatabaseConfig struct {
	ConnectionString string `json:"connection_string"`
}

type LoggerConfig struct {
	LogLevel string `json:"log_level"`
}
