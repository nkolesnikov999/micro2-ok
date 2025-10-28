package order_part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func DeleteOrderParts(ctx context.Context, conn *pgx.Conn, orderUUID uuid.UUID) error {
	query := `DELETE FROM order_parts WHERE order_uuid = $1`
	_, err := conn.Exec(ctx, query, orderUUID)
	return err
}
