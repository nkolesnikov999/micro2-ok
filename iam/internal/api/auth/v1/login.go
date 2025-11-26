package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	login := req.GetLogin()
	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "login cannot be empty")
	}

	password := req.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password cannot be empty")
	}

	sessionUUID, err := a.authService.Login(ctx, login, password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication failed")
	}

	return &authV1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
