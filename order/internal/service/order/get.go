package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, model.ErrOrderGetFailed
	}
	return order, nil
}
