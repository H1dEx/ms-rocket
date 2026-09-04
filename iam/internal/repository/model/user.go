package model

import "time"

type User struct {
	UUID                string               `json:"uuid"`
	Login               string               `json:"login"`
	Email               string               `json:"email"`
	Password            string               `json:"password"`
	NotificationMethods []NotificationMethod `json:"notification_methods"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Endpoint     string `json:"endpoint"`
}
