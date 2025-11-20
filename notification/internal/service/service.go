package service

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/notification/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, orderPaidEvent model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, shipAssembledEvent model.ShipAssembledEvent) error
}
