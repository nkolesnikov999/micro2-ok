package user

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	user, err := s.userRepository.GetUser(ctx, userUUID)
	if err != nil {
		logger.Error(ctx,
			"failed to get user",
			zap.String("userUUID", userUUID),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, model.ErrUserGetFailed
	}
	logger.Debug(ctx,
		"user retrieved successfully",
		zap.String("userUUID", userUUID),
	)
	return user, nil
}
