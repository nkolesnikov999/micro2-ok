package repository

import (
	"context"
	"time"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session model.Session, ttl time.Duration) error
	GetSession(ctx context.Context, sessionUUID string) (model.Session, error)
	AddSessionToUserSet(ctx context.Context, userUUID string, sessionUUID string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User, passwordHash string) error
	GetUser(ctx context.Context, userUUID string) (model.User, error)
	GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (model.User, string, error)
}
