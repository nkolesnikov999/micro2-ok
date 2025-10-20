package part

import (
	"log"
	"sync"

	def "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

func NewRepository() *repository {
	r := &repository{
		parts: make(map[string]*model.Part),
	}

	err := r.initParts(100)
	if err != nil {
		log.Fatalf("failed to initialize parts: %v", err)
	}
	return r
}
