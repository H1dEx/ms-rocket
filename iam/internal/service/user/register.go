package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) RegisterUser(ctx context.Context, login, email, password string, notificationMethods []model.NotificationMethod) (string, error) {
	userUUID := uuid.New().String()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "failed to hash password", zap.Error(err))
		return "", err
	}

	user := model.UserCreate{
		UUID:                userUUID,
		Login:               login,
		Email:               email,
		PasswordHash:        string(hashPassword),
		NotificationMethods: notificationMethods,
	}

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		if !errors.Is(err, model.ErrUserAlreadyExists) {
			logger.Error(ctx, "failed to create user", zap.Error(err))
		}
		return "", err
	}
	return id, nil
}
