package v1

import (
	"context"

	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error) {
	response, err := c.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}
	return response.GetTransactionUuid(), nil
}
