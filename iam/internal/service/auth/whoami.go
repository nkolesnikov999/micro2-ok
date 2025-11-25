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

	// Получаем пользователя из репозитория по userUUID из сессии
	user, err := s.userRepository.GetUser(ctx, session.UserUUID.String())
	if err != nil {
		logger.Error(ctx,
			"failed to get user",
			zap.String("userUUID", session.UserUUID.String()),
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
		zap.String("userUUID", session.UserUUID.String()),
	)

	return session, user, nil
}
