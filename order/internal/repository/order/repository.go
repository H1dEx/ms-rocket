package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/H1dEx/ms-rocket/order/internal/repository"
)

var _ repository.OrderRepository = (*rep)(nil)

type rep struct {
	conn *pgxpool.Pool
}

func NewOrderRepository(conn *pgxpool.Pool) *rep {
	return &rep{
		conn: conn,
	}
}
