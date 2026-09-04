package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/H1dEx/ms-rocket/iam/internal/repository"
)

var _ repository.UserRepository = (*rep)(nil)

type rep struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *rep {
	return &rep{
		conn: conn,
	}
}
