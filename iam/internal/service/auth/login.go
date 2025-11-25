package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) Login(ctx context.Context, login string, password string) (string, error) {
	// Валидация входных данных
	if login == "" {
		logger.Error(ctx, "empty login provided")
		return "", model.ErrInvalidLogin
	}
	if password == "" {
		logger.Error(ctx, "empty password provided")
		return "", model.ErrInvalidPassword
	}

	// Получаем пользователя по login или email вместе с password hash
	user, passwordHash, err := s.userRepository.GetUserByLoginOrEmail(ctx, login)
	if err != nil {
		logger.Error(ctx,
			"failed to get user by login or email",
			zap.String("login", login),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrUserNotFound
		}
		return "", err
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		logger.Error(ctx,
			"invalid password",
			zap.String("login", login),
			zap.Error(err),
		)
		return "", model.ErrInvalidPassword
	}

	// Создаем сессию
	sessionUUID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Сессия действительна 24 часа

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  user.UUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Сохраняем сессию в Redis
	ttl := expiresAt.Sub(now)
	err = s.sessionRepository.CreateSession(ctx, session, ttl)
	if err != nil {
		logger.Error(ctx,
			"failed to create session",
			zap.String("sessionUUID", sessionUUID.String()),
			zap.Error(err),
		)
		return "", err
	}

	// Сохраняем связь сессии с пользователем
	err = s.sessionRepository.AddSessionToUserSet(ctx, user.UUID.String(), sessionUUID.String())
	if err != nil {
		logger.Error(ctx,
			"failed to add session to user set",
			zap.String("userUUID", user.UUID.String()),
			zap.String("sessionUUID", sessionUUID.String()),
			zap.Error(err),
		)
		return "", err
	}

	logger.Debug(ctx,
		"user logged in successfully",
		zap.String("userUUID", user.UUID.String()),
		zap.String("sessionUUID", sessionUUID.String()),
		zap.String("login", login),
	)

	return sessionUUID.String(), nil
}
