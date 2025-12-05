package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	orderMetrics "github.com/nkolesnikov999/micro2-OK/order/internal/metrics"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "order.call_pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID.String()),
			attribute.String("payment.method", paymentMethod),
		),
	)
	defer span.End()

	// Создаем спан для запроса к БД GetOrder
	ctx, dbSpan := tracing.StartSpan(ctx, "db.get_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID.String()),
			attribute.String("operation.name", "pay_order"),
		),
	)
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.End()
		span.RecordError(err)
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
	dbSpan.End()

	if order.Status == "PAID" || order.Status == "CANCELLED" || order.Status == "ASSEMBLED" {
		span.SetAttributes(
			attribute.String("order.status", order.Status),
			attribute.String("payment.method", paymentMethod),
			attribute.String("order.uuid", orderUUID.String()),
		)
		span.SetStatus(codes.Error, "order is not payable")
		logger.Error(ctx,
			"order is not payable",
			zap.String("paymentMethod", paymentMethod),
			zap.String("orderUUID", orderUUID.String()),
			zap.Error(model.ErrOrderNotPayable),
		)
		return "", model.ErrOrderNotPayable
	}

	// Создаем спан для вызова paymentClient
	ctx, clientSpan := tracing.StartSpan(ctx, "grpc.payment.pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID.String()),
			attribute.String("user.uuid", order.UserUUID.String()),
			attribute.String("payment.method", paymentMethod),
			attribute.String("operation.name", "pay_order"),
		),
	)
	txUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.OrderUUID.String(),
		order.UserUUID.String(),
		paymentMethod,
	)
	if err != nil {
		clientSpan.RecordError(err)
		clientSpan.End()
		span.RecordError(err)
		logger.Error(ctx,
			"failed to pay order",
			zap.String("paymentMethod", paymentMethod),
			zap.Any("order", order),
			zap.Error(err),
		)
		return "", model.ErrPaymentFailed
	}
	clientSpan.End()

	order.Status = "PAID"
	order.TransactionUUID = txUUID
	order.PaymentMethod = paymentMethod
	order.UpdatedAt = time.Now()

	// Создаем спан для запроса к БД UpdateOrder
	ctx, updateSpan := tracing.StartSpan(ctx, "db.update_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID.String()),
			attribute.String("order.status", order.Status),
			attribute.String("operation.name", "pay_order"),
		),
	)
	if err := s.orderRepository.UpdateOrder(ctx, orderUUID, order); err != nil {
		updateSpan.RecordError(err)
		updateSpan.End()
		span.RecordError(err)
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrOrderUpdateFailed
	}
	updateSpan.End()

	// Увеличиваем бизнес-метрику выручки на сумму оплаченного заказа.
	// Метрика OrdersRevenueTotal — монотонно возрастающий счетчик общей выручки.
	orderMetrics.OrdersRevenueTotal.Add(ctx, order.TotalPrice)

	err = s.orderPaidProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       orderUUID.String(),
		UserUUID:        order.UserUUID.String(),
		PaymentMethod:   paymentMethod,
		TransactionUUID: txUUID,
	})
	if err != nil {
		span.RecordError(err)
		return "", model.ErrOrderProducerFailed
	}

	span.SetAttributes(
		attribute.String("order.status", order.Status),
		attribute.String("payment.method", paymentMethod),
		attribute.String("order.uuid", orderUUID.String()),
		attribute.String("transactionUUID", txUUID),
	)

	logger.Debug(ctx,
		"order paid successfully",
		zap.String("paymentMethod", paymentMethod),
		zap.Any("order", order),
		zap.String("transactionUUID", txUUID),
	)

	return txUUID, nil
}
