package orderconsumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/nkolesnikov999/micro2-OK/order/internal/converter/kafka"
	"github.com/nkolesnikov999/micro2-OK/order/internal/repository"
	def "github.com/nkolesnikov999/micro2-OK/order/internal/service"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderRepository        repository.OrderRepository
}

func NewService(orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	orderRepository repository.OrderRepository,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderRepository:        orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order assembled consumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
