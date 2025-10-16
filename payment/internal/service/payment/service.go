package payment

import (
	def "github.com/nkolesnikov999/micro2-OK/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
