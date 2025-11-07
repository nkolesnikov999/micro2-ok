package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type HTTPEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type HTTPConfig struct {
	raw HTTPEnvConfig
}

func NewHTTPConfig() (*HTTPConfig, error) {
	var raw HTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &HTTPConfig{raw: raw}, nil
}

func (cfg *HTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *HTTPConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *HTTPConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}
