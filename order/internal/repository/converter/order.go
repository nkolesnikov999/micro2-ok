package converter

import (
	"github.com/google/uuid"
	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/order/internal/repository/model"
)

func ToRepoOrder(order model.Order) repoModel.Order {
	partUuids := make([]string, 0, len(order.PartUuids))
	for _, uuid := range order.PartUuids {
		partUuids = append(partUuids, uuid.String())
	}

	return repoModel.Order{
		OrderUUID:       order.OrderUUID.String(),
		UserUUID:        order.UserUUID.String(),
		PartUuids:       partUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}

func ToModelOrder(order repoModel.Order) model.Order {
	orderUUID, _ := uuid.Parse(order.OrderUUID)
	userUUID, _ := uuid.Parse(order.UserUUID)

	partUuids := make([]uuid.UUID, 0, len(order.PartUuids))
	for _, uuidStr := range order.PartUuids {
		if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
			partUuids = append(partUuids, parsedUUID)
		}
	}

	return model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
