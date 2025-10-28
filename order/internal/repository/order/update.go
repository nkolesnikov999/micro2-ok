package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
	orderpart "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order_part"
)

func (r *repository) UpdateOrder(ctx context.Context, id uuid.UUID, order model.Order) error {
	query := `
		UPDATE orders 
		SET user_uuid = $2, total_price = $3, 
		    transaction_uuid = $4, payment_method = $5, status = $6, updated_at = $7
		WHERE order_uuid = $1`

	repoOrder := repoConverter.ToRepoOrder(order)

	tx, err := r.connDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	result, err := tx.Exec(ctx, query,
		id,
		repoOrder.UserUUID,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
		repoOrder.Status,
		repoOrder.UpdatedAt,
	)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("update orders exec failed and rollback failed: %w", errors.Join(err, rbErr))
		}
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("order not found; rollback failed: %w", rbErr)
		}
		return model.ErrOrderNotFound
	}

	if err := orderpart.UpdateOrderPartsTx(ctx, tx, id, order.PartUuids); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("update order parts failed and rollback failed: %w", errors.Join(err, rbErr))
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
