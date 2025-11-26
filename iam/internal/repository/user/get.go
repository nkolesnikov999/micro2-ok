package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

func (r *repository) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	userUUIDParsed, err := uuid.Parse(userUUID)
	if err != nil {
		return model.User{}, err
	}

	query := `
		SELECT uuid, login, email, password_hash, 
		       notification_methods, created_at, updated_at
		FROM users 
		WHERE uuid = $1`

	rows, err := r.connDB.Query(ctx, query, userUUIDParsed)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	repoUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return repoConverter.ToModelUser(repoUser)
}
