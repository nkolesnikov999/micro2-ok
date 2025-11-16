package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpc "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/mocks"
	repoMocks "github.com/nkolesnikov999/micro2-OK/order/internal/repository/mocks"
	svcMocks "github.com/nkolesnikov999/micro2-OK/order/internal/service/mocks"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository      *repoMocks.OrderRepository
	orderProducerService *svcMocks.OrderProducerService
	paymentClient        *grpc.PaymentClient
	inventoryClient      *grpc.InventoryClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.InitForBenchmark()

	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.orderProducerService = svcMocks.NewOrderProducerService(s.T())
	s.paymentClient = grpc.NewPaymentClient(s.T())
	s.inventoryClient = grpc.NewInventoryClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.orderProducerService,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
