package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCancelById(ctx context.Context, params orderV1.OrderCancelByIdParams) (orderV1.OrderCancelByIdRes, error) {
	err := a.service.OrderCancelById(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
			}, nil
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Conflict error with id %s", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}
	return &orderV1.OrderCancelByIdNoContent{}, nil
}
