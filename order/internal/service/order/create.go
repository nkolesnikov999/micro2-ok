package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	orderMetrics "github.com/nkolesnikov999/micro2-OK/order/internal/metrics"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (model.Order, error) {
	if len(partUUIDs) == 0 {
		logger.Error(ctx,
			"empty part UUIDs",
			zap.String("userUUID", userUUID.String()),
			zap.Any("partUUIDs", partUUIDs),
		)
		return model.Order{}, model.ErrEmptyPartUUIDs
	}

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{Uuids: partUUIDs})
	if err != nil {
		logger.Error(ctx,
			"failed to list parts from inventory",
			zap.String("userUUID", userUUID.String()),
			zap.Any("partUUIDs", partUUIDs),
			zap.Error(err),
		)
		return model.Order{}, model.ErrInventoryUnavailable
	}
	var total float64
	for _, p := range parts {
		total += p.Price
	}

	now := time.Now()
	order := model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: total,
		Status:     "PENDING_PAYMENT",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.orderRepository.CreateOrder(ctx, order, model.PartsFilter{Uuids: partUUIDs}, parts); err != nil {
		logger.Error(ctx,
			"failed to create order",
			zap.Any("order", order),
			zap.Any("partUUIDs", partUUIDs),
			zap.Any("parts", parts),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrOrderAlreadyExists) {
			return model.Order{}, model.ErrOrderAlreadyExists
		}

		var missing *model.PartsNotFoundError
		if errors.As(err, &missing) {
			return model.Order{}, err
		}
		return model.Order{}, model.ErrOrderCreateFailed
	}

	orderMetrics.OrdersTotal.Add(ctx, 1)

	logger.Debug(ctx,
		"order created successfully",
		zap.Any("order", order),
	)

	return order, nil
}
