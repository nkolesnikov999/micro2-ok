package order_part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateOrderParts(ctx context.Context, conn *pgx.Conn, orderUUID uuid.UUID, partUuids []uuid.UUID) error {
	if len(partUuids) == 0 {
		return nil
	}

	query := `INSERT INTO order_parts (order_uuid, part_uuid, quantity) VALUES ($1, $2, $3)`

	for _, partUUID := range partUuids {
		_, err := conn.Exec(ctx, query, orderUUID, partUUID, 1)
		if err != nil {
			return err
		}
	}

	return nil
}
