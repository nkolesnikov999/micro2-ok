package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get order",
			zap.String("orderUUID", orderUUID.String()),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return model.ErrOrderGetFailed
	}

	if order.Status == "PAID" {
		logger.Error(ctx,
			"cannot cancel paid order",
			zap.String("orderUUID", orderUUID.String()),
			zap.Any("order", order),
		)
		return model.ErrCannotCancelPaidOrder
	}

	if order.Status == "PENDING_PAYMENT" {
		order.Status = "CANCELLED"
		order.UpdatedAt = time.Now()
		if err := s.orderRepository.UpdateOrder(ctx, orderUUID, order); err != nil {
			logger.Error(ctx,
				"failed to update order",
				zap.String("orderUUID", orderUUID.String()),
				zap.Any("order", order),
				zap.Error(err),
			)
			if errors.Is(err, model.ErrOrderNotFound) {
				return model.ErrOrderNotFound
			}
			return model.ErrOrderUpdateFailed
		}
	}
	logger.Debug(ctx,
		"order cancelled successfully",
		zap.Any("order", order),
	)
	return nil
}
