package order

import (
	"sync"

	"github.com/H1dEx/ms-rocket/order/internal/repository"
	"github.com/H1dEx/ms-rocket/order/internal/repository/model"
)

var _ repository.OrderRepository = (*rep)(nil)

type rep struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewOrderRepository() *rep {
	return &rep{
		orders: make(map[string]model.Order),
	}
}
