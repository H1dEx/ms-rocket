package v1

import "github.com/H1dEx/ms-rocket/order/internal/service"

type api struct {
	service service.OrderService
}

func NewOrderApi(service service.OrderService) *api {
	return &api{service: service}
}
