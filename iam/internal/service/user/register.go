package user

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

func (s *service) Register(ctx context.Context, login, email, password string) (string, error) {
	// Валидация входных данных
	if login == "" {
		logger.Error(ctx, "empty login provided")
		return "", model.ErrInvalidLogin
	}
	if email == "" {
		logger.Error(ctx, "empty email provided")
		return "", model.ErrInvalidEmail
	}
	if password == "" {
		logger.Error(ctx, "empty password provided")
		return "", model.ErrInvalidPassword
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx,
			"failed to hash password",
			zap.String("email", email),
			zap.Error(err),
		)
		return "", model.ErrPasswordHashFailed
	}

	// Генерируем UUID для пользователя
	userUUID := uuid.New()
	now := time.Now()

	// Создаем модель пользователя
	user := model.User{
		UUID:      userUUID,
		CreatedAt: now,
		UpdatedAt: now,
		Info: model.UserInfo{
			Login:               login,
			Email:               email,
			NotificationMethods: []model.NotificationMethod{},
		},
	}

	// Создаем пользователя в репозитории
	if err := s.userRepository.CreateUser(ctx, user, string(passwordHash)); err != nil {
		logger.Error(ctx,
			"failed to create user",
			zap.String("login", login),
			zap.String("email", email),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return "", model.ErrUserAlreadyExists
		}
		return "", model.ErrUserCreateFailed
	}

	logger.Debug(ctx,
		"user registered successfully",
		zap.String("userUUID", userUUID.String()),
		zap.String("login", login),
		zap.String("email", email),
	)

	return userUUID.String(), nil
}
