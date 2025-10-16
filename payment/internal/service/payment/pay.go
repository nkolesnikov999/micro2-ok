package payment

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod string) (transactionUUID string, err error) {
	method := strings.ToUpper(strings.TrimSpace(paymentMethod))
	switch method {
	case "CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY":
		// ok
	default:
		return "", model.ErrInvalidPaymentMethod
	}

	return uuid.NewString(), nil
}
