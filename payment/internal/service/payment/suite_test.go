package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.InitForBenchmark()

	s.ctx = context.Background()

	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
