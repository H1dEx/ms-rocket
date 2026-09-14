package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/H1dEx/ms-rocket/iam/internal/converter"
	"github.com/H1dEx/ms-rocket/iam/internal/model"
	userV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	if req == nil || req.Info == nil || req.Info.Info == nil {
		return nil, status.Errorf(codes.InvalidArgument, "info is required")
	}

	if req.Info.Info.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, "login is required")
	}

	if req.Info.Info.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	if req.Info.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	userId, err := a.service.RegisterUser(ctx, req.Info.Info.Login, req.Info.Info.Email, req.Info.Password, converter.NotificationMethodsToModel(req.Info.Info.NotificationMethods))
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, err
	}
	return &userV1.RegisterResponse{UserUuid: userId}, nil
}
