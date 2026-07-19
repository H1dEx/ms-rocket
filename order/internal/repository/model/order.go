package model

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusUnknown        OrderStatus = "UNKNOWN"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	// Уникальный идентификатор заказа.
	OrderUUID string `json:"order_uuid"`
	// UUID пользователя.
	UserUUID string `json:"user_uuid"`
	// Список UUID деталей.
	PartUuids []string `json:"part_uuids"`
	// Итоговая стоимость.
	TotalPrice float32 `json:"total_price"`
	// UUID транзакции (если оплачен).
	TransactionUUID *string       `json:"transaction_uuid"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	Status          OrderStatus   `json:"status"`
}
