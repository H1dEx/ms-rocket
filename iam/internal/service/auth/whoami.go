package auth

import (
	"context"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

func (s *service) WhoAmI(ctx context.Context, sessionUUID string) (model.User, model.Session, error) {
	session, err := s.sessionRepository.GetSession(ctx, sessionUUID)
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	user, err := s.userService.GetUser(ctx, session.UserUUID)
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	return user, session, nil
}
