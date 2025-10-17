package grpc

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error)
}
