package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/order/internal/repository/converter"
	orderpart "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order_part"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order, filter model.PartsFilter, parts []model.Part) error {
	if filter.Uuids == nil {
		filter.Uuids = []uuid.UUID{}
	}

	// Проверка наличия всех деталей: сравниваем UUID из фильтра и из parts
	if len(filter.Uuids) > 0 {
		present := make(map[uuid.UUID]struct{}, len(parts))
		for _, p := range parts {
			present[p.Uuid] = struct{}{}
		}
		var missingUUIDs []string
		for _, id := range filter.Uuids {
			if _, ok := present[id]; !ok {
				missingUUIDs = append(missingUUIDs, id.String())
			}
		}
		if len(missingUUIDs) > 0 {
			return &model.PartsNotFoundError{MissingUUIDs: missingUUIDs}
		}
	}
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

	if err := orderpart.CreateOrderParts(ctx, r.connDB, repoOrder.OrderUUID, filter.Uuids); err != nil {
		return err
	}

	return nil
}
