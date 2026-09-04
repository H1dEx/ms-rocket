package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/H1dEx/ms-rocket/iam/internal/converter"
	"github.com/H1dEx/ms-rocket/iam/internal/model"
	authV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	if req == nil || req.SessionUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "session_uuid is required")
	}
	user, session, err := a.service.WhoAmI(ctx, req.SessionUuid)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) || errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.Unauthenticated, "session not found")
		}
		return nil, err
	}

	return &authV1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
