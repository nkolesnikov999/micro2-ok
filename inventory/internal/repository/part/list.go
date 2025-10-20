package part

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/converter"
)

func (r *repository) ListParts(ctx context.Context) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, repoConverter.ToModelPart(*part))
	}
	return parts, nil
}
