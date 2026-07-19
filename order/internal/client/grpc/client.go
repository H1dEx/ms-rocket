package grpc

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (transactionID string, err error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}
