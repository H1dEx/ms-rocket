package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
