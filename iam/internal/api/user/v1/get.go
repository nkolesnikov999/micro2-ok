package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/converter"
	userV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userUUIDStr := req.GetUserUuid()
	if userUUIDStr == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid cannot be empty")
	}

	// Вызываем сервис GetUser, который возвращает User
	user, err := a.userService.GetUser(ctx, userUUIDStr)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Конвертируем User в proto
	protoUser := converter.ToProtoUser(user)

	return &userV1.GetUserResponse{
		User: protoUser,
	}, nil
}
