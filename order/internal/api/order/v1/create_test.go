package v1

import (
	"errors"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	"github.com/brianvoe/gofakeit/v7"
)

func (a *ApiSuite) TestCreateOrderByIDSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		partsUUID = []string{gofakeit.UUID()}
		order     = model.Order{OrderUUID: orderUUID, PartUuids: partsUUID, UserUUID: userUUID, TotalPrice: 10}
		param     = &orderV1.CreateOrderRequest{UserUUID: userUUID, PartUuids: partsUUID}
		res       = &orderV1.CreateOrderResponse{OrderUUID: order.OrderUUID, TotalPrice: order.TotalPrice}
	)

	a.service.On("CreateOrder", a.ctx, userUUID, partsUUID).Return(order, nil).Once()
	response, err := a.api.CreateOrder(a.ctx, param)

	a.NoError(err)
	a.Equal(response, res)
}

func (a *ApiSuite) TestCreateOrderPartsNotFound() {
	var (
		userUUID  = gofakeit.UUID()
		partsUUID = []string{gofakeit.UUID()}
		param     = &orderV1.CreateOrderRequest{UserUUID: userUUID, PartUuids: partsUUID}
		res       = &orderV1.NotFoundError{Code: 404, Message: model.ErrPartsNotFound.Error()}
	)

	a.service.On("CreateOrder", a.ctx, userUUID, partsUUID).Return(model.Order{}, model.ErrPartsNotFound).Once()
	response, err := a.api.CreateOrder(a.ctx, param)

	a.NoError(err)
	a.Equal(response, res)
}

func (a *ApiSuite) TestCreateOrderOrderServerErr() {
	var (
		userUUID  = gofakeit.UUID()
		partsUUID = []string{gofakeit.UUID()}
		param     = &orderV1.CreateOrderRequest{UserUUID: userUUID, PartUuids: partsUUID}
		res       = &orderV1.InternalServerError{Code: 500, Message: "Service error"}
	)

	a.service.On("CreateOrder", a.ctx, userUUID, partsUUID).Return(model.Order{}, errors.New("unknown err")).Once()
	response, err := a.api.CreateOrder(a.ctx, param)

	a.NoError(err)
	a.Equal(response, res)
}
