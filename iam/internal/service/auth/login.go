package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
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

	// Получаем пользователя по login
	user, err := s.getUserByLogin(ctx, login)
	if err != nil {
		logger.Error(ctx,
			"failed to get user by login",
			zap.String("login", login),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrUserNotFound
		}
		return "", err
	}

	// Получаем password hash из репозитория
	passwordHash, err := s.getPasswordHashByUserUUID(ctx, user.UUID.String())
	if err != nil {
		logger.Error(ctx,
			"failed to get password hash",
			zap.String("userUUID", user.UUID.String()),
			zap.Error(err),
		)
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
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Конвертируем в репозиторную модель
	sessionRedisView := repoConverter.ToRepoSession(session)

	// Сохраняем сессию в Redis
	ttl := expiresAt.Sub(now)
	err = s.sessionRepository.CreateSession(ctx, sessionRedisView, ttl)
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

	// Сохраняем userUUID в отдельном ключе Redis для быстрого доступа
	err = s.saveUserUUIDBySessionUUID(ctx, sessionUUID.String(), user.UUID.String())
	if err != nil {
		logger.Error(ctx,
			"failed to save user_uuid for session",
			zap.String("sessionUUID", sessionUUID.String()),
			zap.String("userUUID", user.UUID.String()),
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

// getUserByLogin получает пользователя по login
// TODO: Добавить метод GetUserByLogin в репозиторий пользователей
func (s *service) getUserByLogin(ctx context.Context, login string) (model.User, error) {
	// Временная реализация: возвращаем ошибку
	// В реальной реализации нужно добавить метод GetUserByLogin в репозиторий
	return model.User{}, model.ErrUserNotFound
}

// getPasswordHashByUserUUID получает password hash по userUUID
// TODO: Добавить метод GetPasswordHashByUserUUID в репозиторий пользователей
func (s *service) getPasswordHashByUserUUID(ctx context.Context, userUUID string) (string, error) {
	// Временная реализация: возвращаем ошибку
	// В реальной реализации нужно добавить метод GetPasswordHashByUserUUID в репозиторий
	return "", model.ErrUserNotFound
}

// saveUserUUIDBySessionUUID сохраняет userUUID в отдельном ключе Redis
// Ключ: "iam:session:{sessionUUID}:user"
func (s *service) saveUserUUIDBySessionUUID(ctx context.Context, sessionUUID string, userUUID string) error {
	// TODO: Реализовать сохранение userUUID в Redis
	// Для этого нужно добавить метод в репозиторий сессий
	return nil
}
