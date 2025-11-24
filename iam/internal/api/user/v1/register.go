package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	info := req.GetInfo()
	if info == nil {
		return nil, status.Error(codes.InvalidArgument, "info cannot be nil")
	}

	userInfo := info.GetInfo()
	if userInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "info.info cannot be nil")
	}

	login := userInfo.GetLogin()
	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "login cannot be empty")
	}

	email := userInfo.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}

	password := info.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password cannot be empty")
	}

	userUUID, err := a.userService.Register(ctx, login, email, password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &userV1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
