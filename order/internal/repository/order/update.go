package order

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/order/internal/repository/converter"
)

func (r *rep) UpdateOrder(ctx context.Context, params model.UpdateOrderParam) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[params.OrderUUID]

	if !ok {
		return model.ErrOrderNotFound
	}

	if params.PaymentMethod != nil {
		order.PaymentMethod = converter.PaymentMethodToRepoModel(*params.PaymentMethod)
	}

	if params.Status != nil {
		order.Status = converter.StatusToRepoModel(*params.Status)
	}

	if params.TransactionUUID != nil {
		order.TransactionUUID = params.TransactionUUID
	}

	r.orders[params.OrderUUID] = order
	return nil
}
