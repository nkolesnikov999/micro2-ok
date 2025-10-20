package part

import (
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
	def "github.com/nkolesnikov999/micro2-OK/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
