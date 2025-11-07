package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/mocks"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	partRepository *mocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.InitForBenchmark()

	s.ctx = context.Background()

	s.partRepository = mocks.NewPartRepository(s.T())

	s.service = NewService(
		s.partRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
