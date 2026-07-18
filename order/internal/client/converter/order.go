package converter

import (
	"github.com/H1dEx/ms-rocket/order/internal/model"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	payment_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

func OrderPaymentMethodToModel(payment orderV1.PaymentMethod) model.PaymentMethod {
	switch payment {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func PaymentMethodToOrder(payment model.PaymentMethod) orderV1.PaymentMethod {
	switch payment {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func PaymentMethodToPayment(paymentMethod model.PaymentMethod) payment_v1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

func StatusToOrder(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	default:
		return orderV1.OrderStatusUNKNOWN
	}
}

func OrderToApi(order model.Order) orderV1.OrderDto {
	transaction := orderV1.OptString{Set: false}

	if order.TransactionUUID != nil {
		transaction = orderV1.NewOptString(*order.TransactionUUID)
	}

	return orderV1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transaction,
		PaymentMethod:   orderV1.NewOptPaymentMethod(PaymentMethodToOrder(order.PaymentMethod)),
		Status:          StatusToOrder(order.Status),
	}
}
