package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	GRPC    GRPCConfig
	Tracing TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  loggerCfg,
		GRPC:    grpcCfg,
		Tracing: tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
