package v1

import (
	"context"

	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{StatusCode: 500, Response: orderV1.GenericError{}}
}
