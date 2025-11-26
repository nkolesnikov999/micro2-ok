package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
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
		logger.Error(ctx,
			"failed to register user",
			zap.String("login", login),
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to register user: %v", err))
	}

	return &userV1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
