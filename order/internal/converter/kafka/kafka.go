package kafka

import "github.com/nkolesnikov999/micro2-OK/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
