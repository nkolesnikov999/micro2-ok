package session

import (
	"fmt"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/cache"
)

const (
	cacheKeyPrefix  = "iam:session:"
	userSessionsKey = "iam:user:%s:sessions"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

func (r *repository) getUserSessionsKey(userUUID string) string {
	return fmt.Sprintf(userSessionsKey, userUUID)
}
