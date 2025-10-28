package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderGetFailed
	}

	if order.Status == "PAID" || order.Status == "CANCELLED" {
		return "", model.ErrOrderNotPayable
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
	order.UpdatedAt = time.Now()

	if err := s.orderRepository.UpdateOrder(ctx, orderUUID, order); err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderUpdateFailed
	}

	return txUUID, nil
}
