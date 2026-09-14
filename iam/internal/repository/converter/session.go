package converter

import (
	"time"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	repoModel "github.com/H1dEx/ms-rocket/iam/internal/repository/model"
)

func RedisSessionToModel(redisSession repoModel.SessionRedisView) model.Session {
	createdAt := time.Unix(0, redisSession.CreatedAt).UTC()
	updatedAt := time.Unix(0, redisSession.UpdatedAt).UTC()
	expiresAt := time.Unix(0, redisSession.ExpiresAt).UTC()

	return model.Session{
		UUID:      redisSession.SessionUUID,
		UserUUID:  redisSession.UserUUID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}
}

func SessionToRedisModel(session model.Session) repoModel.SessionRedisView {
	createdAt := session.CreatedAt.UnixNano()
	updatedAt := session.UpdatedAt.UnixNano()
	expiresAt := session.ExpiresAt.UnixNano()
	return repoModel.SessionRedisView{
		SessionUUID: session.UUID,
		UserUUID:    session.UserUUID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		ExpiresAt:   expiresAt,
	}
}
