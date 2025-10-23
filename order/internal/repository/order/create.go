package order

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	// Сначала проверим, существует ли заказ
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uuid = $1)"
	err := r.connDB.QueryRow(ctx, checkQuery, order.OrderUUID.String()).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return model.ErrOrderAlreadyExists
	}

	insertQuery := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, 
		                   transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	repoOrder := repoConverter.ToRepoOrder(order)

	_, err = r.connDB.Exec(ctx, insertQuery,
		repoOrder.OrderUUID,
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

	return nil
}
