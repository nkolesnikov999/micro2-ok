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

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get order",
			zap.String("paymentMethod", paymentMethod),
			zap.String("orderUUID", orderUUID.String()),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderGetFailed
	}

	if order.Status == "PAID" || order.Status == "CANCELLED" || order.Status == "ASSEMBLED" {
		logger.Error(ctx,
			"order is not payable",
			zap.String("paymentMethod", paymentMethod),
			zap.String("orderUUID", orderUUID.String()),
			zap.Error(model.ErrOrderNotPayable),
		)
		return "", model.ErrOrderNotPayable
	}

	txUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.OrderUUID.String(),
		order.UserUUID.String(),
		paymentMethod,
	)
	if err != nil {
		logger.Error(ctx,
			"failed to pay order",
			zap.String("paymentMethod", paymentMethod),
			zap.Any("order", order),
			zap.Error(err),
		)
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

	err = s.orderPaidProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       orderUUID.String(),
		UserUUID:        order.UserUUID.String(),
		PaymentMethod:   paymentMethod,
		TransactionUUID: txUUID,
	})
	if err != nil {
		return "", model.ErrOrderProducerFailed
	}

	logger.Debug(ctx,
		"order paid successfully",
		zap.String("paymentMethod", paymentMethod),
		zap.Any("order", order),
		zap.String("transactionUUID", txUUID),
	)

	return txUUID, nil
}
