package model

import "time"

type NotificationMethod struct {
	ProviderName string
	Endpoint     string
}

type User struct {
	UUID                string
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
