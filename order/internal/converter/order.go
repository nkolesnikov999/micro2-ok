package converter

import (
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	api "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func OrderToAPI(o model.Order) *api.OrderDto {
	partIDs := make([]uuid.UUID, 0, len(o.PartUuids))
	partIDs = append(partIDs, o.PartUuids...)
	dto := &api.OrderDto{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       partIDs,
		TotalPrice:      float32(o.TotalPrice),
		TransactionUUID: api.NewOptNilString(o.TransactionUUID),
		Status:          api.OrderStatus(o.Status),
	}
	if o.PaymentMethod != "" {
		dto.PaymentMethod = api.NewOptPaymentMethod(api.PaymentMethod(o.PaymentMethod))
	}
	return dto
}
