package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/converter"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()
	modelFilter := converter.PartsFilterToModel(filter)

	parts, err := a.inventoryService.ListParts(ctx, modelFilter)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "parts not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	protoParts := converter.PartsToProto(parts)

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
