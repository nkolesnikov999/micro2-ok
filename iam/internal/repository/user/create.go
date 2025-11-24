package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
)

func (r *repository) CreateUser(ctx context.Context, user model.User, passwordHash string) error {
	insertQuery := `
		INSERT INTO users (uuid, login, email, password_hash, 
		                  notification_methods, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	repoUser, err := repoConverter.ToRepoUserWithPasswordHash(user, passwordHash)
	if err != nil {
		return err
	}

	_, err = r.connDB.Exec(ctx, insertQuery,
		repoUser.UUID,
		repoUser.Login,
		repoUser.Email,
		repoUser.PasswordHash,
		repoUser.NotificationMethods,
		repoUser.CreatedAt,
		repoUser.UpdatedAt,
	)
	if err != nil {
		// Проверяем, является ли ошибка нарушением ограничения уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}
