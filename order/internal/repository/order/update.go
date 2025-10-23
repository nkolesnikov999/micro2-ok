package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, id uuid.UUID, order model.Order) error {
	query := `
		UPDATE orders 
		SET user_uuid = $2, part_uuids = $3, total_price = $4, 
		    transaction_uuid = $5, payment_method = $6, status = $7,
		    updated_at = NOW()
		WHERE order_uuid = $1`

	repoOrder := repoConverter.ToRepoOrder(order)

	result, err := r.connDB.Exec(ctx, query,
		id.String(),
		repoOrder.UserUUID,
		repoOrder.PartUuids,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
		repoOrder.Status,
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
