package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/nkolesnikov999/micro2-OK/order/internal/converter"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, err := h.service.GetOrder(ctx, params.OrderUUID)
    if err != nil {
        if errors.Is(err, model.ErrOrderNotFound) {
            return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
        }
        return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
    }

	return converter.OrderToAPI(order), nil
}
