package converter

import (
	"github.com/H1dEx/ms-rocket/notification/internal/model"
	v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

func PaymentMethodToModel(method v1.PaymentMethod) model.PaymentMethod {
	switch method {
	case v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
