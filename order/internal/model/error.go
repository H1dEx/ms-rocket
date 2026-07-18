package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrNotPendingStatus = errors.New("order status is not pending")
var ErrPartsNotFound = errors.New("parts not found")