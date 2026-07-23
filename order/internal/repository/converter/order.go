package converter

import (
	"github.com/H1dEx/ms-rocket/order/internal/model"
	repoModel "github.com/H1dEx/ms-rocket/order/internal/repository/model"
)

func OrderToModel(o repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.OrderUUID,
		PartUuids:       o.PartUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   PaymentMethodToModel(o.PaymentMethod),
		Status:          StatusToModel(o.Status),
	}
}

func PaymentMethodToModel(method *string) model.PaymentMethod {
	if method == nil {
		return model.PaymentMethodUnknown
	}
	switch *method {
	case string(model.PaymentMethodCard):
		return model.PaymentMethodCard
	case string(model.PaymentMethodSBP):
		return model.PaymentMethodSBP
	case string(model.PaymentMethodCreditCard):
		return model.PaymentMethodCreditCard
	case string(model.PaymentMethodInvestorMoney):
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func StatusToModel(status string) model.OrderStatus {
	switch status {
	case string(model.OrderStatusPendingPayment):
		return model.OrderStatusPendingPayment
	case string(model.OrderStatusPaid):
		return model.OrderStatusPaid
	case string(model.OrderStatusCancelled):
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusUnknown
	}
}
