package repository

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderUUID, userUUId string, partUuids []string, price float32) error
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
	UpdateOrder(ctx context.Context, params model.UpdateOrderParam) error
}
