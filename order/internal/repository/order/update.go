package order

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (r *rep) UpdateOrder(ctx context.Context, params model.UpdateOrderParam) error {
	var query strings.Builder

	query.WriteString("UPDATE orders SET ")

	sets := make([]string, 0, 4)
	values := make([]interface{}, 0, 4)

	if params.PaymentMethod != nil {
		sets = append(sets, fmt.Sprintf("payment_method = $%d", len(values)+1))
		values = append(values, *params.PaymentMethod)
	}

	if params.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", len(values)+1))
		values = append(values, *params.Status)
	}

	if params.TransactionUUID != nil {
		sets = append(sets, fmt.Sprintf("transaction_uuid = $%d", len(values)+1))
		values = append(values, *params.TransactionUUID)
	}

	values = append(values, params.OrderUUID)

	if len(sets) == 0 {
		return model.ErrInvalidUpdateOrderParams
	}

	query.WriteString(strings.Join(sets, ", "))

	fmt.Fprintf(&query, " WHERE order_uuid = $%d", len(values))
	res, err := r.conn.Exec(ctx, query.String(), values...)
	if err != nil {
		log.Printf("error updating order: %v", err)
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
