package config

import "github.com/spf13/viper"

const (
	configDir = "./config"
)

func ReadConfigYmlFile() *TelegramConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("err") // TODO: Handle error through logger
		}
	}

	tc := NewTelegramConfig(
		viper.GetString("telegram.api_token"),
	)

	return tc
}
