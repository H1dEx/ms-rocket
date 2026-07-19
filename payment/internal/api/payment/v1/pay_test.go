package v1

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"

	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

func (a *ApiSuite) TestPayOrderSuccess() {
	uuid := gofakeit.UUID()
	a.service.On("PayOrder", a.ctx).Return(uuid, nil).Once()

	res, err := a.api.PayOrder(a.ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	a.NoError(err)
	a.Equal(&paymentV1.PayOrderResponse{
		TransactionUuid: uuid,
	}, res)
}

func (a *ApiSuite) TestPayOrderError() {
	ErrTest := errors.New("Test error expected")
	a.service.On("PayOrder", a.ctx).Return("", ErrTest).Once()

	res, err := a.api.PayOrder(a.ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	a.Error(err)
	a.ErrorIs(err, ErrTest)
	a.Empty(res)
}
