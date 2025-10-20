package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrder(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, uuid uuid.UUID, order model.Order) error
}
