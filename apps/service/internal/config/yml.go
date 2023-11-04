package config

import "github.com/spf13/viper"

const (
	configFilePath = "/etc/kitadrivebot-service.yaml"
)

func ReadConfig() (*Config, error) {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
