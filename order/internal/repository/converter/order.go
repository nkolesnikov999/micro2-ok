package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/order/internal/repository/model"
)

func ToRepoOrder(order model.Order) repoModel.Order {
	var transactionUUID uuid.UUID
	if order.TransactionUUID != "" {
		if parsed, err := uuid.Parse(order.TransactionUUID); err == nil {
			transactionUUID = parsed
		}
	}

	return repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func ToModelOrder(order repoModel.Order, partUuids []uuid.UUID) model.Order {
	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       partUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID.String(),
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
