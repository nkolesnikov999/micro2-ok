package repository

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

// PartRepository определяет интерфейс для работы с частями в репозитории
type PartRepository interface {
	// GetPart возвращает часть по UUID
	GetPart(ctx context.Context, uuid string) (model.Part, error)

	// ListParts возвращает список частей с возможностью фильтрации
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)

	// CreatePart создает новую часть
	InitParts(ctx context.Context, parts []repoModel.Part) error
}
