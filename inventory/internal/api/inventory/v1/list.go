package v1

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/converter"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()

	// Конвертируем proto фильтр в model фильтр
	modelFilter := converter.PartsFilterToModel(filter)

	// Получение частей через сервисный слой
	parts, err := a.inventoryService.ListParts(ctx, modelFilter)
	if err != nil {
		return nil, err
	}

	// Конвертируем model части в proto части
	protoParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		protoPart := converter.PartToProto(part)
		protoParts = append(protoParts, protoPart)
	}

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
