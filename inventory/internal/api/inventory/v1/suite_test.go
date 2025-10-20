package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	inventoryService *mocks.PartService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryService = mocks.NewPartService(s.T())

	s.api = NewAPI(
		s.inventoryService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
