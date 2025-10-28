package order_part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func UpdateOrderParts(ctx context.Context, conn *pgx.Conn, orderUUID uuid.UUID, partUuids []uuid.UUID) error {
	// Сначала удаляем все существующие части заказа
	deleteQuery := `DELETE FROM order_parts WHERE order_uuid = $1`
	_, err := conn.Exec(ctx, deleteQuery, orderUUID)
	if err != nil {
		return err
	}

	// Затем добавляем новые части
	return CreateOrderParts(ctx, conn, orderUUID, partUuids)
}
