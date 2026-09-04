package v1

import (
	"github.com/H1dEx/ms-rocket/iam/internal/service"
	authV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthApi(srv service.AuthService) authV1.AuthServiceServer {
	return &api{service: srv}
}
