package payment

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod string) (transactionUUID string, err error) {
	method := strings.ToUpper(strings.TrimSpace(paymentMethod))
	switch method {
	case "CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY":
		// ok
	default:
		logger.Error(ctx,
			"invalid payment method",
			zap.String("paymentMethod", paymentMethod),
		)
		return "", model.ErrInvalidPaymentMethod
	}

	logger.Debug(ctx,
		"payment method validated successfully",
		zap.String("paymentMethod", paymentMethod),
	)

	return uuid.NewString(), nil
}
