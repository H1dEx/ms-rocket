package session

import (
	"fmt"

	"github.com/H1dEx/ms-rocket/iam/internal/repository"
	"github.com/H1dEx/ms-rocket/platform/pkg/cache"
)

const (
	cacheHashKeyPrefix = "iam:session:"
	cacheSetKeyPrefix  = "iam:user:"
)

type rep struct {
	cache cache.RedisClient
}

var _ repository.SessionRepository = (*rep)(nil)

func NewRepository(cache cache.RedisClient) *rep {
	return &rep{
		cache: cache,
	}
}

func (r *rep) getHashCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheHashKeyPrefix, uuid)
}

func (r *rep) getSetHashCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheSetKeyPrefix, uuid)
}
