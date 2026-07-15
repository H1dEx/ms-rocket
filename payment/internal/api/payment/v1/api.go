package v1

import (
	"github.com/H1dEx/ms-rocket/payment/internal/service"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewApi(service service.PaymentService) *api {
	return &api{
		paymentService: service,
	}
}
