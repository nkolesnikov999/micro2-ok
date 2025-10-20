package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/nkolesnikov999/micro2-OK/order/internal/converter"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if req == nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
	}

	paymentMethod := converter.PaymentMethodToModel(req.PaymentMethod)
	tx, err := h.service.PayOrder(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, model.ErrOrderNotPayable):
			return &orderV1.ConflictError{Code: http.StatusConflict, Message: err.Error()}, nil
		case errors.Is(err, model.ErrPaymentFailed):
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}

	return &orderV1.PayOrderResponse{TransactionUUID: tx}, nil
}
