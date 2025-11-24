package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/converter"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	sessionUUIDStr := req.GetSessionUuid()
	if sessionUUIDStr == "" {
		return nil, status.Error(codes.InvalidArgument, "session_uuid cannot be empty")
	}

	// Вызываем сервис Whoami, который возвращает Session и User
	session, user, err := a.authService.Whoami(ctx, sessionUUIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "session not found or expired")
	}

	// Конвертируем Session и User в proto
	protoSession := converter.ToProtoSession(session)
	protoUser := converter.ToProtoUser(user)

	return &authV1.WhoamiResponse{
		Session: protoSession,
		User:    protoUser,
	}, nil
}
