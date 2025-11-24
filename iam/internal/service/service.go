package service

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
}

type UserService interface {
	GetUser(ctx context.Context, uuid string) (model.User, error)
	Register(ctx context.Context, login string, email string, password string) (string, error)
}
