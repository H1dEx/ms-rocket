package order

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (r *rep) CreateOrder(ctx context.Context, orderUUID, userUUID string, partUuids []string, price float32) error {
	res, err := r.conn.Exec(ctx, "INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, status) VALUES ($1, $2, $3, $4, $5)", orderUUID, userUUID, partUuids, price, string(model.OrderStatusPendingPayment))
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotCreated
	}

	return nil
}
