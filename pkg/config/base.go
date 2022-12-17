package config

type TelegramConfig struct {
	Token string
}

func NewTelegramConfig(token string) *TelegramConfig {
	return &TelegramConfig{
		Token: token,
	}
}
