package config

import (
	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	Config() *sarama.Config
	GroupID() string
}

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}
