package model

type Order struct {
	OrderUUID       string   `db:"order_uuid"`
	UserUUID        string   `db:"user_uuid"`
	PartUuids       []string `db:"part_uuids"`
	TotalPrice      float64  `db:"total_price"`
	TransactionUUID string   `db:"transaction_uuid"`
	PaymentMethod   string   `db:"payment_method"`
	Status          string   `db:"status"`
}
