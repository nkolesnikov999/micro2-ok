package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (model.Order, error) {
	if len(partUUIDs) == 0 {
		return model.Order{}, model.ErrEmptyPartUUIDs
	}

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{Uuids: partUUIDs})
	if err != nil {
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
		if errors.Is(err, model.ErrOrderAlreadyExists) {
			return model.Order{}, model.ErrOrderAlreadyExists
		}

		var missing *model.PartsNotFoundError
		if errors.As(err, &missing) {
			return model.Order{}, err
		}
		return model.Order{}, model.ErrOrderCreateFailed
	}

	return order, nil
}
