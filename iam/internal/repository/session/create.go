package session

import (
	"context"
	"fmt"
	"time"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	"github.com/H1dEx/ms-rocket/iam/internal/repository/converter"
)

func (r *rep) CreateSession(ctx context.Context, session model.Session, ttl time.Duration) error {
	cacheKey := r.getHashCacheKey(session.UUID)

	redisSession := converter.SessionToRedisModel(session)

	err := r.cache.HashSet(ctx, cacheKey, redisSession)
	if err != nil {
		return fmt.Errorf("save session hash: %w", err)
	}

	err = r.cache.Expire(ctx, cacheKey, ttl)
	if err != nil {
		return fmt.Errorf("expire session hash: %w", err)
	}

	return r.AddSessionToUser(ctx, session)
}
