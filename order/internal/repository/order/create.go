package order

import (
	"context"

	repoModel "github.com/H1dEx/ms-rocket/order/internal/repository/model"
)

func (r *rep) CreateOrder(ctx context.Context, orderUUID, userUUId string, partUuids []string, price float32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := repoModel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUId,
		PartUuids:  partUuids,
		TotalPrice: price,
		Status:     repoModel.OrderStatusPendingPayment,
	}

	r.orders[orderUUID] = order

	return nil
}
