package model

import "errors"

var (
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrOrderNotFound           = errors.New("order not found")
	ErrEmptyPartUUIDs          = errors.New("part_uuids must not be empty")
	ErrPartsNotFound           = errors.New("one or more parts not found")
	ErrOrderAlreadyPaid        = errors.New("order already paid")
	ErrCannotPayCancelledOrder = errors.New("cannot pay cancelled order")
	ErrCannotCancelPaidOrder   = errors.New("order already paid and cannot be cancelled")

	// Service-level failure categories
	ErrInventoryUnavailable = errors.New("inventory service unavailable")
	ErrPaymentFailed        = errors.New("payment failed")
	ErrOrderCreateFailed    = errors.New("order create failed")
	ErrOrderUpdateFailed    = errors.New("order update failed")
	ErrOrderGetFailed       = errors.New("order get failed")
)
