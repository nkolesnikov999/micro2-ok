package repository

import (
	"context"
	"time"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session repoModel.SessionRedisView, ttl time.Duration) error
	GetSession(ctx context.Context, sessionUUID string) (model.Session, error)
	AddSessionToUserSet(ctx context.Context, userUUID string, sessionUUID string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User, passwordHash string) error
	GetUser(ctx context.Context, userUUID string) (model.User, error)
}
