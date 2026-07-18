package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order, err := a.service.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return &orderV1.NotFoundError{Code: 404, Message: err.Error()}, nil
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Getting order error"),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}
	return &orderV1.CreateOrderResponse{OrderUUID: order.OrderUUID, TotalPrice: order.TotalPrice}, nil
}
