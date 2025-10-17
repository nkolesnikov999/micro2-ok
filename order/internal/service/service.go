package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

type OrderService interface {
	// CreateOrder validates parts via Inventory, calculates total, and persists the order.
	// Returns the created domain order.
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (model.Order, error)

	// GetOrder returns the domain order by its UUID.
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)

	// PayOrder processes payment for the order and returns the transaction UUID.
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (string, error)

	// CancelOrder cancels the order if not paid.
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
}
