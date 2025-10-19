package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nkolesnikov999/micro2-OK/order/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	orderService *mocks.OrderService

	api *orderHandler
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.orderService = mocks.NewOrderService(s.T())

	s.api = NewHandler(
		s.orderService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
