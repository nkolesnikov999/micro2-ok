package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req == nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
	}

	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: model.ErrEmptyPartUUIDs.Error()}, nil
	}

	order, err := h.service.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyPartUUIDs):
			return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
		case errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, model.ErrInventoryUnavailable):
			return &orderV1.ServiceUnavailableError{Code: http.StatusServiceUnavailable, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: float32(order.TotalPrice),
	}, nil
}
