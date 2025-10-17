package part

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(*repoPart), nil
}
