package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type PaymentGRPCEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type PaymentGRPCConfig struct {
	raw PaymentGRPCEnvConfig
}

func NewPaymentGRPCConfig() (*PaymentGRPCConfig, error) {
	var raw PaymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &PaymentGRPCConfig{raw: raw}, nil
}

func (cfg *PaymentGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
