package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := h.service.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, model.ErrCannotCancelPaidOrder):
			return &orderV1.ConflictError{Code: http.StatusConflict, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}
	return nil, &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderV1.GenericError{},
	}
}
