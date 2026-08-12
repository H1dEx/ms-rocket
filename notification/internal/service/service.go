package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/notification/internal/model"
)

type OrderPaidService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledService interface {
	RunConsumer(ctx context.Context) error
}
type TelegramService interface {
	SendPaidNotification(ctx context.Context, data model.OrderPaidEvent) error
	SendAssembledNotification(ctx context.Context, data model.ShipAssembledEvent) error
}
