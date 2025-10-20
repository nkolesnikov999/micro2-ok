package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id.String()]
	if !exists {
		return model.Order{}, model.ErrOrderNotFound
	}
	return repoConverter.ToModelOrder(*order), nil
}
