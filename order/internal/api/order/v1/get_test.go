package v1

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/H1dEx/ms-rocket/order/internal/client/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *ApiSuite) TestGetOrderByIDSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		order     = model.Order{OrderUUID: orderUUID}
		param     = orderV1.GetOrderByIDParams{OrderUUID: orderUUID}
		res       = &orderV1.GetOrderResponse{Order: converter.OrderToApi(order)}
	)

	a.service.On("GetOrderByID", a.ctx, orderUUID).Return(order, nil).Once()
	response, err := a.api.GetOrderByID(a.ctx, param)

	a.NoError(err)
	a.Equal(response, res)
}

func (a *ApiSuite) TestGetOrderByNotFoundErr() {
	var (
		orderUUID = gofakeit.UUID()
		param     = orderV1.GetOrderByIDParams{OrderUUID: orderUUID}
		resErr    = &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with id %s is not found", orderUUID),
		}
	)

	a.service.On("GetOrderByID", a.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()
	response, err := a.api.GetOrderByID(a.ctx, param)

	a.NoError(err)
	a.Equal(response, resErr)
}

func (a *ApiSuite) TestGetOrderByServerErr() {
	var (
		orderUUID = gofakeit.UUID()
		param     = orderV1.GetOrderByIDParams{OrderUUID: orderUUID}
		resErr    = &orderV1.InternalServerError{Code: 500, Message: "Service error"}
	)

	a.service.On("GetOrderByID", a.ctx, orderUUID).Return(model.Order{}, errors.New("error")).Once()
	response, err := a.api.GetOrderByID(a.ctx, param)

	a.NoError(err)
	a.Equal(response, resErr)
}
