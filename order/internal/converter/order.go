package converter

import (
	"github.com/H1dEx/ms-rocket/order/internal/model"
	v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
	payment_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(method model.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

func OrderPaidEventToModel(event model.OrderPaidEvent) v1.OrderPaid {
	return v1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   PaymentMethodToProto(event.PaymentMethod),
		TransactionUuid: event.TransactionUUID,
	}
}

func ShipAssembledEventToModel(event *v1.ShipAssembled) model.ShipAssembledEvent {
	return model.ShipAssembledEvent{
		EventUUID:        event.EventUuid,
		OrderUUID:        event.OrderUuid,
		UserUUID:         event.UserUuid,
		BuildTimeSeconds: event.BuildTimeSeconds,
	}
}
