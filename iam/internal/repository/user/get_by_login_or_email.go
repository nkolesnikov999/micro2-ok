package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

func (r *repository) GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (model.User, string, error) {
	if loginOrEmail == "" {
		return model.User{}, "", model.ErrUserNotFound
	}

	query := `
		SELECT uuid, login, email, password_hash, 
		       notification_methods, created_at, updated_at
		FROM users 
		WHERE login = $1 OR email = $1`

	rows, err := r.connDB.Query(ctx, query, loginOrEmail)
	if err != nil {
		return model.User{}, "", err
	}
	defer rows.Close()

	repoUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, "", model.ErrUserNotFound
		}
		return model.User{}, "", err
	}

	user, err := repoConverter.ToModelUser(repoUser)
	if err != nil {
		return model.User{}, "", err
	}

	return user, repoUser.PasswordHash, nil
}
