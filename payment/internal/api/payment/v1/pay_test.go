package v1

import (
	"errors"

	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
	payment_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v7"
)

func (a *ApiSuite) TestPayOrderSuccess() {
	uuid := gofakeit.UUID()
	a.service.On("PayOrder", a.ctx).Return(uuid, nil).Once()

	res, err := a.api.PayOrder(a.ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	a.NoError(err)
	a.Equal(&paymentV1.PayOrderResponse{
		TransactionUuid: uuid,
	}, res)
}
func (a *ApiSuite) TestPayOrderError() {
	var ErrTest = errors.New("Test error expected")
	a.service.On("PayOrder", a.ctx).Return("", ErrTest).Once()

	res, err := a.api.PayOrder(a.ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	a.Error(err)
	a.ErrorIs(err, ErrTest)
	a.Empty(res)
}
