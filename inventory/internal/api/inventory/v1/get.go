package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/converter"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	partUUID := req.GetUuid()

	if _, err := uuid.Parse(partUUID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	part, err := a.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}

	protoPart := converter.PartToProto(part)
	return &inventoryV1.GetPartResponse{Part: protoPart}, nil
}
