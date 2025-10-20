package repository

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)

	ListParts(ctx context.Context) ([]model.Part, error)
}
