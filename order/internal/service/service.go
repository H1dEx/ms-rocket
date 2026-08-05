package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error)
	GetOrderByID(ctx context.Context, orderUUID string) (model.Order, error)
	OrderCancelByID(ctx context.Context, orderUUID string) error
	PayOrderByID(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error)
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}