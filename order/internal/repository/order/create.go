package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	insertQuery := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, 
		                   transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	repoOrder := repoConverter.ToRepoOrder(order)

	_, err := r.connDB.Exec(ctx, insertQuery,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		repoOrder.PartUuids,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
		repoOrder.Status,
	)
	if err != nil {
		// Проверяем, является ли ошибка нарушением ограничения уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrOrderAlreadyExists
		}
		return err
	}

	return nil
}
