package grpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
	commonV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/common/v1"
)

const (
	// SessionUUIDMetadataKey ключ для передачи UUID сессии в gRPC metadata
	SessionUUIDMetadataKey = "session-uuid"
)

type contextKey string

const (
	// userContextKey ключ для хранения пользователя в контексте
	userContextKey contextKey = "user"
	// sessionUUIDContextKey ключ для хранения session UUID в контексте
	sessionUUIDContextKey contextKey = "session-uuid"
)

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = authV1.AuthServiceClient

// AuthInterceptor interceptor для аутентификации gRPC запросов
type AuthInterceptor struct {
	iamClient IAMClient
}

// NewAuthInterceptor создает новый interceptor аутентификации
func NewAuthInterceptor(iamClient IAMClient) *AuthInterceptor {
	return &AuthInterceptor{
		iamClient: iamClient,
	}
}

// Unary возвращает unary server interceptor для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate выполняет аутентификацию и добавляет пользователя в контекст
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Извлекаем metadata из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Warn(ctx, "[AuthInterceptor] missing metadata")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Логируем все metadata для отладки
	keys := make([]string, 0, len(md))
	for k := range md {
		keys = append(keys, k)
	}
	logger.Info(ctx, "[AuthInterceptor] All metadata keys",
		zap.Strings("keys", keys),
	)
	for k, v := range md {
		logger.Info(ctx, "[AuthInterceptor] Metadata entry",
			zap.String("key", k),
			zap.Strings("values", v),
		)
	}

	// Получаем session UUID из metadata
	sessionUUIDs := md.Get(SessionUUIDMetadataKey)
	if len(sessionUUIDs) == 0 {
		logger.Warn(ctx, "[AuthInterceptor] missing session-uuid in metadata",
			zap.String("expected_key", SessionUUIDMetadataKey),
		)
		return nil, status.Error(codes.Unauthenticated, "missing session-uuid in metadata")
	}

	sessionUUID := sessionUUIDs[0]
	if sessionUUID == "" {
		logger.Warn(ctx, "[AuthInterceptor] empty session-uuid")
		return nil, status.Error(codes.Unauthenticated, "empty session-uuid")
	}

	logger.Info(ctx, "[AuthInterceptor] Found session UUID",
		zap.String("session_uuid", sessionUUID),
	)

	// Валидируем сессию через IAM сервис
	whoamiRes, err := i.iamClient.Whoami(ctx, &authV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		logger.Warn(ctx, "[AuthInterceptor] invalid session",
			zap.String("session_uuid", sessionUUID),
			zap.Error(err),
		)
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid session: %v", err))
	}

	logger.Info(ctx, "[AuthInterceptor] Authentication successful",
		zap.String("user_uuid", whoamiRes.User.Uuid),
		zap.String("session_uuid", sessionUUID),
	)

	// Добавляем пользователя и session UUID в контекст
	authCtx := context.WithValue(ctx, userContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, sessionUUIDContextKey, sessionUUID)
	return authCtx, nil
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	user, ok := ctx.Value(userContextKey).(*commonV1.User)
	return user, ok
}

// GetUserContextKey возвращает ключ контекста для пользователя
func GetUserContextKey() contextKey {
	return userContextKey
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}

// AddSessionUUIDToContext добавляет session UUID в контекст
func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionUUIDContextKey, sessionUUID)
}

// ForwardSessionUUIDToGRPC добавляет session UUID из контекста в исходящие gRPC metadata
func ForwardSessionUUIDToGRPC(ctx context.Context) context.Context {
	sessionUUID, ok := GetSessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionUUIDMetadataKey, sessionUUID)
}
