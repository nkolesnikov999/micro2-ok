package v1

import (
	"context"

	grpcAuth "github.com/nkolesnikov999/micro2-OK/platform/pkg/middleware/grpc"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error) {
	// Передаем session UUID в gRPC metadata для аутентификации
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	var methodEnum paymentV1.PaymentMethod
	switch paymentMethod {
	case "CARD":
		methodEnum = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		methodEnum = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		methodEnum = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY":
		methodEnum = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		methodEnum = paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	response, err := c.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: methodEnum,
	})
	if err != nil {
		return "", err
	}
	return response.GetTransactionUuid(), nil
}
