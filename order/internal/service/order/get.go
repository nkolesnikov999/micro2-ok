package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get order",
			zap.String("orderUUID", orderUUID.String()),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, model.ErrOrderGetFailed
	}
	logger.Debug(ctx,
		"order retrieved successfully",
		zap.Any("order", order),
	)
	return order, nil
}
