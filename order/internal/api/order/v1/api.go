package v1

import (
	"github.com/H1dEx/ms-rocket/order/internal/service"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

type api struct {
	service service.OrderService
}

func NewOrderAPI(svc service.OrderService) orderV1.Handler {
	return &api{service: svc}
}
