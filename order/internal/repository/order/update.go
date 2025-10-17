package order

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, uuid string, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[uuid]; !exists {
		return model.ErrOrderNotFound
	}

	repoOrder := repoConverter.OrderToRepoModel(order)
	r.orders[uuid] = &repoOrder
	return nil
}
