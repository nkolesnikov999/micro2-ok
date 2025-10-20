package v1

import (
	"github.com/nkolesnikov999/micro2-OK/order/internal/service"
)

type orderHandler struct {
	service service.OrderService
}

func NewHandler(service service.OrderService) *orderHandler {
	return &orderHandler{service: service}
}
