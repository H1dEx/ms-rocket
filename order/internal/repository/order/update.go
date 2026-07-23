package order

import (
	"context"
	"log"
	"strings"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (r *rep) UpdateOrder(ctx context.Context, params model.UpdateOrderParam) error {
	var query strings.Builder

	query.WriteString("UPDATE orders SET ")

	if params.PaymentMethod != nil {
		query.WriteString("payment_method = $1, ")
	}

	if params.Status != nil {
		query.WriteString("status = $2, ")
	}

	if params.TransactionUUID != nil {
		query.WriteString("transaction_uuid = $3 ")
	}

	query.WriteString("WHERE order_uuid = $4")

	res, err := r.conn.Exec(ctx, query.String(), params.PaymentMethod, params.Status, params.TransactionUUID, params.OrderUUID)
	if err != nil {
		log.Printf("error updating order: %v", err)
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
