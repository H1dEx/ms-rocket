package v1

import (
	"github.com/H1dEx/ms-rocket/iam/internal/service"
	userV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserAPI(svc service.UserService) userV1.UserServiceServer {
	return &api{service: svc}
}
