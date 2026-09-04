package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, login, email, password string, notificationMethods []model.NotificationMethod) (string, error)
	GetUser(ctx context.Context, userUUID string) (model.User, error)
	CheckPassword(ctx context.Context, login, password string) (model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	WhoAmI(ctx context.Context, sessionUUID string) (model.User, model.Session, error)
}
