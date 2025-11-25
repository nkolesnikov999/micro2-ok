package session

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

// GetSession реализует интерфейс SessionRepository
func (r *repository) GetSession(ctx context.Context, sessionUUID string) (model.Session, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Session{}, model.ErrSessionNotFound
		}
		return model.Session{}, err
	}

	if len(values) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return model.Session{}, err
	}

	return repoConverter.ToModelSession(sessionRedisView)
}
