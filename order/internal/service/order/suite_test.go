package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpc "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/mocks"
	"github.com/nkolesnikov999/micro2-OK/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *mocks.OrderRepository
	paymentClient   *grpc.PaymentClient
	inventoryClient *grpc.InventoryClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.paymentClient = grpc.NewPaymentClient(s.T())
	s.inventoryClient = grpc.NewInventoryClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
