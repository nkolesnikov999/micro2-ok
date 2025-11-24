package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	// Валидация входных данных
	if sessionUUID == "" {
		logger.Error(ctx, "empty session_uuid provided")
		return model.Session{}, model.User{}, model.ErrSessionNotFound
	}

	// Получаем сессию из репозитория
	session, err := s.sessionRepository.GetSession(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get session",
			zap.String("sessionUUID", sessionUUID),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrSessionNotFound) {
			return model.Session{}, model.User{}, model.ErrSessionNotFound
		}
		return model.Session{}, model.User{}, err
	}

	// Получаем userUUID из сессии
	// Предполагаем, что userUUID хранится в отдельном ключе Redis "iam:session:{uuid}:user"
	userUUID, err := s.getUserUUIDBySessionUUID(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get user_uuid from session",
			zap.String("sessionUUID", sessionUUID),
			zap.Error(err),
		)
		return model.Session{}, model.User{}, err
	}

	// Получаем пользователя из репозитория
	user, err := s.userRepository.GetUser(ctx, userUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get user",
			zap.String("userUUID", userUUID),
			zap.String("sessionUUID", sessionUUID),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserNotFound) {
			return model.Session{}, model.User{}, model.ErrUserNotFound
		}
		return model.Session{}, model.User{}, err
	}

	logger.Debug(ctx,
		"whoami retrieved successfully",
		zap.String("sessionUUID", sessionUUID),
		zap.String("userUUID", userUUID),
	)

	return session, user, nil
}

// getUserUUIDBySessionUUID получает userUUID из отдельного ключа Redis
// Ключ: "iam:session:{sessionUUID}:user"
func (s *service) getUserUUIDBySessionUUID(ctx context.Context, sessionUUID string) (string, error) {
	// TODO: Реализовать получение userUUID из Redis
	// Для этого нужно добавить метод в репозиторий сессий
	// или использовать отдельный ключ Redis
	// Временная реализация: возвращаем ошибку
	return "", model.ErrSessionNotFound
}
