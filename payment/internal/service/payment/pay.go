package payment

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod string) (transactionUUID string, err error) {
	ctx, span := tracing.StartSpan(ctx, "payment.call_pay_order",
		trace.WithAttributes(
			attribute.String("payment.method", paymentMethod),
		),
	)
	defer span.End()

	method := strings.ToUpper(strings.TrimSpace(paymentMethod))
	switch method {
	case "CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY":
		// ok
	default:
		span.RecordError(model.ErrInvalidPaymentMethod)
		span.SetStatus(codes.Error, "invalid payment method")
		logger.Error(ctx,
			"invalid payment method",
			zap.String("paymentMethod", paymentMethod),
		)
		return "", model.ErrInvalidPaymentMethod
	}

	transactionUUID = uuid.NewString()
	span.SetAttributes(
		attribute.String("payment.method", method),
		attribute.String("transaction.uuid", transactionUUID),
	)

	logger.Debug(ctx,
		"payment method validated successfully",
		zap.String("paymentMethod", paymentMethod),
		zap.String("transactionUUID", transactionUUID),
	)

	return transactionUUID, nil
}
