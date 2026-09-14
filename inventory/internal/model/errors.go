package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")
	ErrUserNotFound = errors.New("user not found")
)
