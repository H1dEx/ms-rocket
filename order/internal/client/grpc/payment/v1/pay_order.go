package v1

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/client/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (transactionID string, err error) {
	res, err := c.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{OrderUuid: orderUuid, UserUuid: userUuid, PaymentMethod: converter.PaymentMethodToPayment(paymentMethod)})
	if err != nil {
		return "", err
	}

	return res.GetTransactionUuid(), nil
}
