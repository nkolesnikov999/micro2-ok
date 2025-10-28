package order_part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func UpdateOrderParts(ctx context.Context, conn *pgx.Conn, orderUUID uuid.UUID, partUuids []uuid.UUID) (err error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("rollback tx: %w", errors.Join(err, rbErr))
			}
			return
		}
		if cmErr := tx.Commit(ctx); cmErr != nil {
			err = fmt.Errorf("commit tx: %w", cmErr)
		}
	}()

	// 1) Удаляем все существующие части заказа
	_, err = tx.Exec(ctx, `DELETE FROM order_parts WHERE order_uuid = $1`, orderUUID)
	if err != nil {
		return fmt.Errorf("delete order_parts: %w", err)
	}

	// 2) Если новый список не пуст, вставляем его батчем
	if len(partUuids) > 0 {
		_, err = tx.Exec(ctx, `
INSERT INTO order_parts (order_uuid, part_uuid, quantity)
SELECT $1::uuid, UNNEST($2::uuid[]), 1
`, orderUUID, partUuids)
		if err != nil {
			return fmt.Errorf("insert order_parts: %w", err)
		}
	}

	return nil
}
