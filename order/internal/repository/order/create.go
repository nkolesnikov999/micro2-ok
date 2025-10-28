package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
	orderpart "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order_part"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	insertQuery := `
		INSERT INTO orders (order_uuid, user_uuid, total_price, 
		                   transaction_uuid, payment_method, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	repoOrder := repoConverter.ToRepoOrder(order)

	_, err := r.connDB.Exec(ctx, insertQuery,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
		repoOrder.Status,
		repoOrder.CreatedAt,
		repoOrder.UpdatedAt,
	)
	if err != nil {
		// Проверяем, является ли ошибка нарушением ограничения уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrOrderAlreadyExists
		}
		return err
	}

	if err := orderpart.CreateOrderParts(ctx, r.connDB, repoOrder.OrderUUID, order.PartUuids); err != nil {
		return err
	}

	return nil
}
