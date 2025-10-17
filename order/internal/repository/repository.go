package repository

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrder(ctx context.Context, uuid string) (model.Order, error)
	UpdateOrder(ctx context.Context, uuid string, order model.Order) error
}
