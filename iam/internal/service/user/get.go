package user

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	user, err := s.userRepository.GetUserByUUID(ctx, userUUID)
	if err != nil {
		if !errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "failed to get user", zap.Error(err))
		}
		return model.User{}, err
	}
	return user, nil
}
