package session

import (
	"context"

	"github.com/gomodule/redigo/redis"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	"github.com/H1dEx/ms-rocket/iam/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/iam/internal/repository/model"
)

func (r *rep) GetSession(ctx context.Context, sessionUUID string) (model.Session, error) {
	cacheKey := r.getHashCacheKey(sessionUUID)
	session, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		return model.Session{}, err
	}

	if len(session) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var redisSession repoModel.SessionRedisView
	err = redis.ScanStruct(session, &redisSession)
	if err != nil {
		return model.Session{}, err
	}

	return converter.RedisSessionToModel(redisSession), nil
}
