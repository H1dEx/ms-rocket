package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/H1dEx/ms-rocket/order/internal/client/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrderById(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByIdParams) (orderV1.PayOrderByIdRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Order uuid is empty",
		}, nil
	}

	paymentMethod := converter.OrderPaymentMethodToModel(req.PaymentMethod)
	if paymentMethod == model.PaymentMethodUnknown {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Unknown payment method %s", req.PaymentMethod),
		}, nil
	}

	transactionUUID, err := a.service.PayOrderById(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
			}, nil
		}
		if errors.Is(err, model.ErrNotPendingStatus) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Conflict error with id %s", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil 
}
