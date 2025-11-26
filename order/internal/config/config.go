package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/nkolesnikov999/micro2-OK/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	HTTP                   HTTPConfig
	Postgres               PostgresConfig
	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	InventoryGRPC          InventoryGRPCConfig
	PaymentGRPC            PaymentGRPCConfig
	IAMGRPC                IAMGRPCConfig
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

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}
	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}
	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}
	iamGRPCCfg, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		HTTP:                   httpCfg,
		Postgres:               postgresCfg,
		InventoryGRPC:          inventoryGRPCCfg,
		PaymentGRPC:            paymentGRPCCfg,
		IAMGRPC:                iamGRPCCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
