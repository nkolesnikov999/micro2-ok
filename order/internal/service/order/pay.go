package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderGetFailed
	}

	if order.Status == "PAID" {
		return "", model.ErrOrderAlreadyPaid
	}
	if order.Status == "CANCELLED" {
		return "", model.ErrCannotPayCancelledOrder
	}

	txUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.OrderUUID.String(),
		order.UserUUID.String(),
		paymentMethod,
	)
	if err != nil {
		return "", model.ErrPaymentFailed
	}

	order.Status = "PAID"
	order.TransactionUUID = txUUID
	order.PaymentMethod = paymentMethod

	if err := s.orderRepository.UpdateOrder(ctx, orderUUID.String(), order); err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderUpdateFailed
	}

	return txUUID, nil
}
