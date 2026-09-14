package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	authV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	if req == nil || req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	sessionUUID, err := a.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, err
	}

	return &authV1.LoginResponse{SessionUuid: sessionUUID}, nil
}
