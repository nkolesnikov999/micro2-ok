package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type sessionEnvConfig struct {
	TTL time.Duration `env:"SESSION_TTL,required"`
}

type sessionConfig struct {
	raw sessionEnvConfig
}

func NewSessionConfig() (*sessionConfig, error) {
	var raw sessionEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &sessionConfig{raw: raw}, nil
}

func (cfg *sessionConfig) TTL() time.Duration {
	return cfg.raw.TTL
}
