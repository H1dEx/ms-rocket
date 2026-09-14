package converter

import (
	"time"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	sessionV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/common/v1"
)

func SessionToModel(session *sessionV1.Session) model.Session {
	if session == nil {
		return model.Session{}
	}
	var (
		createdAt time.Time
		updatedAt time.Time
		expiresAt time.Time
	)

	if session.CreatedAt != nil {
		createdAt = session.CreatedAt.AsTime()
	}
	if session.UpdatedAt != nil {
		createdAt = session.UpdatedAt.AsTime()
	}
	if session.ExpiresAt != nil {
		expiresAt = session.ExpiresAt.AsTime()
	}
	return model.Session{
		UUID:      session.Uuid,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}
}
