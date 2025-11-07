package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		logger.Error(ctx,
			"failed to get part",
			zap.String("uuid", uuid),
			zap.Error(err),
		)
		return model.Part{}, err
	}

	logger.Debug(ctx,
		"part retrieved successfully",
		zap.Any("part", part),
	)

	return part, nil
}
