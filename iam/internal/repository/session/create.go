package session

import (
	"context"
	"time"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
)

func (r *repository) CreateSession(ctx context.Context, session model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(session.UUID.String())

	redisView := repoConverter.ToRepoSession(session)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
