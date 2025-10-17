package order

import (
	"context"
	"errors"

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

	found := make(map[uuid.UUID]struct{}, len(parts))
	var total float64
	for _, p := range parts {
		found[p.Uuid] = struct{}{}
		total += p.Price
	}

	for _, id := range partUUIDs {
		if _, ok := found[id]; !ok {
			return model.Order{}, model.ErrPartsNotFound
		}
	}

	order := model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: total,
		Status:     "PENDING_PAYMENT",
	}

	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		if errors.Is(err, model.ErrOrderAlreadyExists) {
			return model.Order{}, model.ErrOrderAlreadyExists
		}
		return model.Order{}, model.ErrOrderCreateFailed
	}

	return order, nil
}
