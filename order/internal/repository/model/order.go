package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID `db:"order_uuid"`
	UserUUID        uuid.UUID `db:"user_uuid"`
	PartUuids       []string  `db:"part_uuids"`
	TotalPrice      float64   `db:"total_price"`
	TransactionUUID uuid.UUID `db:"transaction_uuid"`
	PaymentMethod   string    `db:"payment_method"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
