package app

import (
	paymentApi "github.com/H1dEx/ms-rocket/payment/internal/api/payment/v1"
	"github.com/H1dEx/ms-rocket/payment/internal/service"
	paymentService "github.com/H1dEx/ms-rocket/payment/internal/service/payment"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) PaymentV1API() paymentV1.PaymentServiceServer {
	if c.paymentV1API == nil {
		c.paymentV1API = paymentApi.NewAPI(c.PaymentService())
	}
	return c.paymentV1API
}

func (c *diContainer) PaymentService() service.PaymentService {
	if c.paymentService == nil {
		c.paymentService = paymentService.NewService()
	}
	return c.paymentService
}
