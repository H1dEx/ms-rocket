package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/H1dEx/ms-rocket/iam/internal/converter"
	"github.com/H1dEx/ms-rocket/iam/internal/model"
	userV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	if req == nil || req.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	if _, err := uuid.Parse(req.UserUuid); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_uuid")
	}

	user, err := a.service.GetUser(ctx, req.UserUuid)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with UUID %s not found", req.UserUuid)
		}
		return nil, err
	}
	return &userV1.GetUserResponse{User: converter.UserToProto(user)}, nil
}
