package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/order/internal/repository/model"
)

func (r *repository) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	query := `
		SELECT order_uuid, user_uuid, part_uuids, total_price, 
		       transaction_uuid, payment_method, status, created_at, updated_at
		FROM orders 
		WHERE order_uuid = $1`

	rows, err := r.connDB.Query(ctx, query, id)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	repoOrder, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}

	return repoConverter.ToModelOrder(repoOrder), nil
}
