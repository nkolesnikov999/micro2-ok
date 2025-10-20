package order

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return model.ErrOrderGetFailed
	}

	if order.Status == "PAID" {
		return model.ErrCannotCancelPaidOrder
	}

	if order.Status == "PENDING_PAYMENT" {
		order.Status = "CANCELLED"
		if err := s.orderRepository.UpdateOrder(ctx, orderUUID, order); err != nil {
			if errors.Is(err, model.ErrOrderNotFound) {
				return model.ErrOrderNotFound
			}
			return model.ErrOrderUpdateFailed
		}
	}
	return nil
}
