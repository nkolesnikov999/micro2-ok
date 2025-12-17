package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamauth "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *iamauth.LoginRequest) (*iamauth.LoginResponse, error) {
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

	return &iamauth.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
