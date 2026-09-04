package session

import (
	"context"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
)

func (r *rep) AddSessionToUser(ctx context.Context, session model.Session) error {
	cacheKey := r.getSetHashCacheKey(session.UserUUID)
	return r.cache.SAdd(ctx, cacheKey, session.UUID)
}
