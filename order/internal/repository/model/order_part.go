package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderPart struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
}
