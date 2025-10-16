package part

import (
	"sync"

	def "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

// Repository представляет in-memory репозиторий для работы с частями
type repository struct {
	mu    sync.RWMutex
	parts map[string]model.Part
}

// NewRepository создает новый экземпляр репозитория
func NewRepository() *repository {
	return &repository{
		parts: make(map[string]model.Part),
	}
}
