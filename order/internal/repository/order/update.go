package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, id uuid.UUID, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := id.String()
	if _, exists := r.orders[key]; !exists {
		return model.ErrOrderNotFound
	}

	repoOrder := repoConverter.ToRepoOrder(order)
	r.orders[key] = &repoOrder
	return nil
}
