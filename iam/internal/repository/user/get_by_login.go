package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	"github.com/H1dEx/ms-rocket/iam/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/iam/internal/repository/model"
)

func (r *rep) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var (
		UUID                    string
		email                   string
		notificationMethodsJSON []byte
		notificationMethods     []repoModel.NotificationMethod
		passwordHash            string
		createdAt               time.Time
		updatedAt               time.Time
	)

	err := r.conn.QueryRow(ctx, "SELECT uuid, email, notification_methods, password, created_at, updated_at FROM users WHERE login = $1 LIMIT 1", login).Scan(&UUID, &email, &notificationMethodsJSON, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}
	err = json.Unmarshal(notificationMethodsJSON, &notificationMethods)
	if err != nil {
		return model.User{}, err
	}

	repoUser := repoModel.User{
		UUID:                UUID,
		Login:               login,
		Email:               email,
		NotificationMethods: notificationMethods,
		Password:            passwordHash,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}

	return converter.UserToModel(repoUser), nil
}
