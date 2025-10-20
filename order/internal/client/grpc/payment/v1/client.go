package v1

import (
	def "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	paymentClient paymentV1.PaymentServiceClient
}

func NewClient(paymentClient paymentV1.PaymentServiceClient) *client {
	return &client{
		paymentClient: paymentClient,
	}
}
