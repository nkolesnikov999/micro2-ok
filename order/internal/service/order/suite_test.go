package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"

	grpc "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/mocks"
	orderMetrics "github.com/nkolesnikov999/micro2-OK/order/internal/metrics"
	repoMocks "github.com/nkolesnikov999/micro2-OK/order/internal/repository/mocks"
	svcMocks "github.com/nkolesnikov999/micro2-OK/order/internal/service/mocks"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository      *repoMocks.OrderRepository
	orderProducerService *svcMocks.OrderPaidProducerService
	paymentClient        *grpc.PaymentClient
	inventoryClient      *grpc.InventoryClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.InitForBenchmark()

	// Initialize no-op metrics provider for tests
	noopProvider := metric.NewMeterProvider()
	otel.SetMeterProvider(noopProvider)

	// Initialize order service metrics
	_ = orderMetrics.InitMetrics("order-service-test")

	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.orderProducerService = svcMocks.NewOrderPaidProducerService(s.T())
	s.paymentClient = grpc.NewPaymentClient(s.T())
	s.inventoryClient = grpc.NewInventoryClient(s.T())

	// By default, allow producing OrderPaidRecorded without error
	s.orderProducerService.
		On("ProduceOrderPaid", mock.Anything, mock.Anything).
		Return(nil).
		Maybe()

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
