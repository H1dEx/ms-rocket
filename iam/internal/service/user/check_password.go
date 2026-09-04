package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

func (s *service) CheckPassword(ctx context.Context, login, password string) (model.User, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, model.ErrInvalidCredentials
		}
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return model.User{}, model.ErrInvalidCredentials
	}

	user.PasswordHash = ""

	return user, nil
}
