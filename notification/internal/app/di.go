package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/nkolesnikov999/micro2-OK/notification/internal/client/http"
	telegramClient "github.com/nkolesnikov999/micro2-OK/notification/internal/client/http/telegram"
	"github.com/nkolesnikov999/micro2-OK/notification/internal/config"
	kafkaConverter "github.com/nkolesnikov999/micro2-OK/notification/internal/converter/kafka"
	kafkaDecoder "github.com/nkolesnikov999/micro2-OK/notification/internal/converter/kafka/decoder"
	"github.com/nkolesnikov999/micro2-OK/notification/internal/service"
	orderPaidConsumer "github.com/nkolesnikov999/micro2-OK/notification/internal/service/consumer/order_paid_consumer"
	shipAssembledConsumer "github.com/nkolesnikov999/micro2-OK/notification/internal/service/consumer/ship_assembled_consumer"
	telegramService "github.com/nkolesnikov999/micro2-OK/notification/internal/service/telegram"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	wrappedKafka "github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/consumer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService
	telegramService               service.TelegramService

	orderPaidConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      wrappedKafka.Consumer
	orderAssembledDecoder       kafkaConverter.OrderAssembledDecoder

	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(ctx),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = shipAssembledConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.TelegramService(ctx),
		)
	}
	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
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

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
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

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = kafkaDecoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}
