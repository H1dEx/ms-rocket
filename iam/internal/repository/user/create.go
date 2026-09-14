package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

func (r *rep) CreateUser(ctx context.Context, user model.UserCreate) (string, error) {
	createdAt := time.Now()
	if user.NotificationMethods == nil {
		user.NotificationMethods = []model.NotificationMethod{}
	}
	notificationMethod, err := json.Marshal(user.NotificationMethods)
	if err != nil {
		return "", err
	}
	notificationMethodJSON := string(notificationMethod)

	res, err := r.conn.Exec(ctx, "INSERT INTO users (uuid, login, email, notification_methods, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", user.UUID, user.Login, user.Email, notificationMethodJSON, user.PasswordHash, createdAt, createdAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// "23505" is the standard Postgres code for unique_violation
			if pgErr.Code == "23505" {
				return "", model.ErrUserAlreadyExists
			}
		}
		return "", err
	}

	if res.RowsAffected() == 0 {
		return "", model.ErrUserNotCreated
	}

	return user.UUID, nil
}
