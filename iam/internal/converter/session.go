package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	sessionV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/common/v1"
)

func SessionToProto(session model.Session) *sessionV1.Session {
	createdAt := timestamppb.New(session.CreatedAt)
	updatedAt := timestamppb.New(session.UpdatedAt)
	expiresAt := timestamppb.New(session.ExpiresAt)
	return &sessionV1.Session{
		Uuid:      session.UUID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}
}
