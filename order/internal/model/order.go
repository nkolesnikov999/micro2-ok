package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
