package model

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotCreated     = errors.New("user not created")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrSessionNotFound = errors.New("session not found")
)
