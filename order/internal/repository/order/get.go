package order

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/order/internal/repository/converter"
)

func(r *rep) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]

	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}