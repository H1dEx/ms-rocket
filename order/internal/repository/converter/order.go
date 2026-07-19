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

func PaymentMethodToModel(method repoModel.PaymentMethod) model.PaymentMethod {
	switch method {
	case repoModel.PaymentMethodCard:
		return model.PaymentMethodCard
	case repoModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repoModel.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case repoModel.PaymentMethodInvestorMoney:
		return model.PaymentMethod(repoModel.PaymentMethodInvestorMoney)
	default:
		return model.PaymentMethodUnknown
	}
}

func PaymentMethodToRepoModel(method model.PaymentMethod) repoModel.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return repoModel.PaymentMethodCard
	case model.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return repoModel.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return repoModel.PaymentMethod(repoModel.PaymentMethodInvestorMoney)
	default:
		return repoModel.PaymentMethodUnknown
	}
}

func StatusToModel(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment
	case repoModel.OrderStatusPaid:
		return model.OrderStatusPaid
	case repoModel.OrderStatusCancelled:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusUnknown
	}
}

func StatusToRepoModel(status model.OrderStatus) repoModel.OrderStatus {
	switch status {
	case model.OrderStatusPendingPayment:
		return repoModel.OrderStatusPendingPayment
	case model.OrderStatusPaid:
		return repoModel.OrderStatusPaid
	case model.OrderStatusCancelled:
		return repoModel.OrderStatusCancelled
	default:
		return repoModel.OrderStatusUnknown
	}
}
