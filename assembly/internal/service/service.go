package service

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/assembly/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
