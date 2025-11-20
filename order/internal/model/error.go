package model

import (
	"errors"
	"fmt"
)

var (
	ErrOrderAlreadyExists    = errors.New("order already exists")
	ErrOrderNotFound         = errors.New("order not found")
	ErrEmptyPartUUIDs        = errors.New("part_uuids must not be empty")
	ErrPartsNotFound         = errors.New("one or more parts not found")
	ErrOrderNotPayable       = errors.New("order cannot be paid")
	ErrCannotCancelPaidOrder = errors.New("order already paid and cannot be cancelled")

	// Service-level failure categories
	ErrInventoryUnavailable = errors.New("inventory service unavailable")
	ErrPaymentFailed        = errors.New("payment failed")
	ErrOrderCreateFailed    = errors.New("order create failed")
	ErrOrderUpdateFailed    = errors.New("order update failed")
	ErrOrderGetFailed       = errors.New("order get failed")
	ErrOrderProducerFailed  = errors.New("order producer failed")
)

// PartsNotFoundError содержит информацию об отсутствующих деталях
type PartsNotFoundError struct {
	MissingUUIDs []string
}

func (e *PartsNotFoundError) Error() string {
	return fmt.Sprintf("one or more parts not found: missing UUIDs: %v", e.MissingUUIDs)
}

func (e *PartsNotFoundError) Unwrap() error {
	return ErrPartsNotFound
}
