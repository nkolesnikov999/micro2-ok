package ship_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/consumer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	if err := s.telegramService.SendOrderAssembledNotification(ctx, event); err != nil {
		return err
	}

	return nil
}
