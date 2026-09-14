package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userService.CheckPassword(ctx, login, password)
	if err != nil {
		return "", err
	}

	createdAt := time.Now().UTC()

	session := model.Session{
		UUID:      uuid.New().String(),
		UserUUID:  user.UUID,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		ExpiresAt: createdAt.Add(s.sessionTTL),
	}
	err = s.sessionRepository.CreateSession(ctx, session, s.sessionTTL)
	if err != nil {
		return "", err
	}
	return session.UUID, nil
}
