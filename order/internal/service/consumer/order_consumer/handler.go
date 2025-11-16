package orderconsumer

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
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

	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to parse order UUID",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err))
		return err
	}

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err))
		if errors.Is(err, model.ErrOrderNotFound) {
			// Логируем, но не возвращаем ошибку, чтобы не зацикливать обработку
			logger.Warn(ctx, "Order not found, skipping status update",
				zap.String("order_uuid", event.OrderUUID))
			return nil
		}
		return err
	}

	order.Status = "ASSEMBLED"
	order.UpdatedAt = time.Now()

	err = s.orderRepository.UpdateOrder(ctx, orderUUID, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec))

	return nil
}
