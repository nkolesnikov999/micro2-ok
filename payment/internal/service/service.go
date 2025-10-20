package service

import "context"

type PaymentService interface {
	PayOrder(ctx context.Context, paymentMethod string) (transactionUUID string, err error)
}
