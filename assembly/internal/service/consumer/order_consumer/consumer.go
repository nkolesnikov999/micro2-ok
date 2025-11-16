package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/nkolesnikov999/micro2-OK/assembly/internal/converter/kafka"
	def "github.com/nkolesnikov999/micro2-OK/assembly/internal/service"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

var _ def.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
