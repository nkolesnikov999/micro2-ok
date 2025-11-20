package env

import (
	"github.com/caarlos0/env/v11"
)

type telegramBotEnvConfig struct {
	Token string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type telegramBotConfig struct {
	raw telegramBotEnvConfig
}

func NewTelegramBotConfig() (*telegramBotConfig, error) {
	var raw telegramBotEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramBotConfig{raw: raw}, nil
}

func (cfg *telegramBotConfig) Token() string {
	return cfg.raw.Token
}
