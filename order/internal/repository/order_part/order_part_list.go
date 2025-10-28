package order_part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func ListOrderParts(ctx context.Context, conn *pgx.Conn, orderUUID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT part_uuid FROM order_parts WHERE order_uuid = $1`
	rows, err := conn.Query(ctx, query, orderUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	partUuids := make([]uuid.UUID, 0)
	for rows.Next() {
		var partUUID uuid.UUID
		if err := rows.Scan(&partUUID); err != nil {
			return nil, err
		}
		partUuids = append(partUuids, partUUID)
	}

	return partUuids, rows.Err()
}
