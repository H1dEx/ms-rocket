package auth

import (
	"time"

	"github.com/H1dEx/ms-rocket/iam/internal/repository"
	serviceModel "github.com/H1dEx/ms-rocket/iam/internal/service"
)

var _ serviceModel.AuthService = (*service)(nil)

type service struct {
	sessionRepository repository.SessionRepository
	userService       serviceModel.UserService
	sessionTTL        time.Duration
}

func NewService(sessionRepository repository.SessionRepository, userService serviceModel.UserService, sessionTTL time.Duration) *service {
	return &service{
		sessionRepository: sessionRepository,
		userService:       userService,
		sessionTTL:        sessionTTL,
	}
}
