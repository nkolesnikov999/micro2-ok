package order_consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/assembly/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/consumer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
	)

	// start measuring build time from the moment we received the event
	start := time.Now()
	// wait 10 seconds (respecting cancellation)
	select {
	case <-time.After(10 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	elapsed := time.Since(start)

	if err := s.shipAssembledProducer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(elapsed / time.Second),
	}); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
