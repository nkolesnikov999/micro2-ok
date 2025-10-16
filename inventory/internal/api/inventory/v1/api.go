package v1

import (
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/service"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
