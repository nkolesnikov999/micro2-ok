package order

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return model.Order{}, model.ErrOrderNotFound
	}
	return repoConverter.OrderToModel(*order), nil
}
