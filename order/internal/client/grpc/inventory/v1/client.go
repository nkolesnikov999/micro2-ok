package v1

import (
	def "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	inventoryClient inventoryV1.InventoryServiceClient
}

func NewClient(inventoryClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		inventoryClient: inventoryClient,
	}
}
