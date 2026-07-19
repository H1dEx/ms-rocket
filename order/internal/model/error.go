package model

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrNotPendingStatus = errors.New("order status is not pending")
	ErrPartsNotFound    = errors.New("parts not found")
)
