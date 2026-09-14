package repository

import (
	"context"
	"time"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.UserCreate) (string, error)
	GetUserByUUID(ctx context.Context, userUUID string) (model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session model.Session, ttl time.Duration) error
	GetSession(ctx context.Context, sessionUUID string) (model.Session, error)
	AddSessionToUser(ctx context.Context, session model.Session) error
}
