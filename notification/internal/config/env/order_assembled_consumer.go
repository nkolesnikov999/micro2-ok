package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	raw orderAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{raw: raw}, nil
}

func (cfg *orderAssembledConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderAssembledConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderAssembledConsumerConfig) Config() *sarama.Config {
	return newSaramaConsumerConfig()
}
