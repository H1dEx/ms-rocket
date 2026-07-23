package model

type PaymentMethod string

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
	TransactionUUID *string `json:"transaction_uuid"`
	PaymentMethod   *string `json:"payment_method"`
	Status          string  `json:"status"`
}
