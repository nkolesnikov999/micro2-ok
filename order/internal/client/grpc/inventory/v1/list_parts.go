package v1

import (
	"context"

	clientConverter "github.com/nkolesnikov999/micro2-OK/order/internal/client/converter"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.ToProtoPartsFilter(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.ToModelPartList(parts.Parts), nil
}
