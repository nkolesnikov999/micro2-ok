package v1

import (
	"github.com/nkolesnikov999/micro2-OK/payment/internal/service"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
