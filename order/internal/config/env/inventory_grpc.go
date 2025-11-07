package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type InventoryGRPCEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type InventoryGRPCConfig struct {
	raw InventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*InventoryGRPCConfig, error) {
	var raw InventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &InventoryGRPCConfig{raw: raw}, nil
}

func (cfg *InventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
