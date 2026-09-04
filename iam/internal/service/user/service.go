package user

import (
	"github.com/H1dEx/ms-rocket/iam/internal/repository"
	def "github.com/H1dEx/ms-rocket/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
