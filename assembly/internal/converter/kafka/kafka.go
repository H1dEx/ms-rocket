package kafka

import "github.com/H1dEx/ms-rocket/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
