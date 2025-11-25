package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamGRPCConfig struct {
	raw iamGRPCEnvConfig
}

func NewGRPCConfig() (*iamGRPCConfig, error) {
	var raw iamGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamGRPCConfig{raw: raw}, nil
}

func (cfg *iamGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
