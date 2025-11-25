package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type IAMGRPCEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type IAMGRPCConfig struct {
	raw IAMGRPCEnvConfig
}

func NewIAMGRPCConfig() (*IAMGRPCConfig, error) {
	var raw IAMGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &IAMGRPCConfig{raw: raw}, nil
}

func (cfg *IAMGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
