package converter

import (
	"github.com/H1dEx/ms-rocket/assembly/internal/model"
	events_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
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

func OrderPaidEventToModel(event *events_v1.OrderPaid) model.OrderPaidEvent {
	return model.OrderPaidEvent{
		EventUUID:       event.GetEventUuid(),
		OrderUUID:       event.GetOrderUuid(),
		UserUUID:        event.GetUserUuid(),
		PaymentMethod:   PaymentMethodToModel(event.GetPaymentMethod()),
		TransactionUUID: event.GetTransactionUuid(),
	}
}

func ShipAssembledEventToEvent(event model.ShipAssembledEvent) events_v1.ShipAssembled {
	return events_v1.ShipAssembled{
		EventUuid:        event.EventUUID,
		OrderUuid:        event.OrderUUID,
		UserUuid:         event.UserUUID,
		BuildTimeSeconds: event.BuildTimeSeconds,
	}
}
