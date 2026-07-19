package payment

import def "github.com/H1dEx/ms-rocket/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
