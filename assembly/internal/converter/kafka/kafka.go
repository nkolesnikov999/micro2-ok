package kafka

import "github.com/nkolesnikov999/micro2-OK/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
