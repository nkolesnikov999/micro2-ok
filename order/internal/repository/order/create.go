package order

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderUUID := order.OrderUUID.String()
	if _, exists := r.orders[orderUUID]; exists {
		return model.ErrOrderAlreadyExists
	}

	repoOrder := repoConverter.OrderToRepoModel(order)
	r.orders[orderUUID] = &repoOrder
	return nil
}
