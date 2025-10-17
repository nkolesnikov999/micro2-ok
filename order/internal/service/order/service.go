package order

import (
	"github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc"
	"github.com/nkolesnikov999/micro2-OK/order/internal/repository"
	def "github.com/nkolesnikov999/micro2-OK/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
