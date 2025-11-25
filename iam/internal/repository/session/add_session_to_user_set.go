package session

import (
	"context"
)

func (r *repository) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	setKey := r.getUserSessionsKey(userUUID)
	if err := r.cache.SAdd(ctx, setKey, sessionUUID); err != nil {
		return err
	}
	return nil
}
