package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/H1dEx/ms-rocket/order/internal/client/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByID(ctx context.Context, params orderV1.GetOrderByIDParams) (orderV1.GetOrderByIDRes, error) {
	order, err := a.service.GetOrderByID(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}

	return &orderV1.GetOrderResponse{Order: converter.OrderToApi(order)}, nil
}
