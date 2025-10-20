package v1

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/converter"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()
	modelFilter := converter.PartsFilterToModel(filter)

	parts, err := a.inventoryService.ListParts(ctx, modelFilter)
	if err != nil {
		return nil, err
	}

	protoParts := converter.PartsToProto(parts)

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
