package order

import (
	"sync"

	def "github.com/nkolesnikov999/micro2-OK/order/internal/repository"
	"github.com/nkolesnikov999/micro2-OK/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*model.Order),
	}
}
