package kafka

import "github.com/H1dEx/ms-rocket/order/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
