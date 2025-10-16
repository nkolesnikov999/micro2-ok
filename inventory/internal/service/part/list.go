package part

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
