package order

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/order/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/order/internal/repository/model"
)

func (r *rep) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	var (
		userUUID        string
		partUuids       []string
		totalPrice      float32
		transactionUUID *string
		paymentMethod   *string
		status          string
	)

	err := r.conn.QueryRow(ctx, "SELECT user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status FROM orders WHERE order_uuid = $1 LIMIT 1", orderUUID).Scan(&userUUID, &partUuids, &totalPrice, &transactionUUID, &paymentMethod, &status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}

	repoOrder := repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}

	return converter.OrderToModel(repoOrder), nil
}
