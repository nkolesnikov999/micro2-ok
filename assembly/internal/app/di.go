package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/nkolesnikov999/micro2-OK/assembly/internal/config"
	kafkaConverter "github.com/nkolesnikov999/micro2-OK/assembly/internal/converter/kafka"
	kafkaDecoder "github.com/nkolesnikov999/micro2-OK/assembly/internal/converter/kafka/decoder"
	"github.com/nkolesnikov999/micro2-OK/assembly/internal/service"
	orderConsumer "github.com/nkolesnikov999/micro2-OK/assembly/internal/service/consumer/order_consumer"
	shipProducer "github.com/nkolesnikov999/micro2-OK/assembly/internal/service/producer/ship_producer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	wrappedKafka "github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/producer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumerService     service.OrderPaidConsumerService
	shipAssembledProducerService service.ShipAssembledProducerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	syncProducer          sarama.SyncProducer
	shipAssembledProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) ShipAssembledProducerService() service.ShipAssembledProducerService {
	if d.shipAssembledProducerService == nil {
		d.shipAssembledProducerService = shipProducer.NewService(d.ShipAssembledProducer())
	}

	return d.shipAssembledProducerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = kafkaDecoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ShipAssembledProducer() wrappedKafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.shipAssembledProducer
}
