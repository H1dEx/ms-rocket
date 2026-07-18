package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error)
	GetOrderByID(ctx context.Context, orderUUID string) (model.Order, error)
	OrderCancelById(ctx context.Context, orderUUID string) error
	PayOrderById(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error)
}
