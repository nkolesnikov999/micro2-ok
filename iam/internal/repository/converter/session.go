package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

// ToRepoSession конвертирует доменную модель сессии в репозиторную модель для Redis.
func ToRepoSession(session model.Session) repoModel.SessionRedisView {
	var updatedAtNs *int64
	// Устанавливаем UpdatedAt только если он отличается от CreatedAt
	if !session.UpdatedAt.Equal(session.CreatedAt) {
		updatedAtNsValue := session.UpdatedAt.UnixNano()
		updatedAtNs = &updatedAtNsValue
	}

	return repoModel.SessionRedisView{
		UUID:        session.UUID.String(),
		UserUUID:    session.UserUUID.String(),
		CreatedAtNs: session.CreatedAt.UnixNano(),
		UpdatedAtNs: updatedAtNs,
		ExpiresAtNs: session.ExpiresAt.UnixNano(),
	}
}

// ToModelSession конвертирует репозиторную модель сессии из Redis в доменную модель.
func ToModelSession(redisView repoModel.SessionRedisView) (model.Session, error) {
	uuidVal, err := uuid.Parse(redisView.UUID)
	if err != nil {
		return model.Session{}, err
	}

	userUUID, err := uuid.Parse(redisView.UserUUID)
	if err != nil {
		return model.Session{}, err
	}

	createdAt := time.Unix(0, redisView.CreatedAtNs)

	var updatedAt time.Time
	if redisView.UpdatedAtNs != nil {
		updatedAt = time.Unix(0, *redisView.UpdatedAtNs)
	} else {
		// Если UpdatedAt не установлен, используем CreatedAt
		updatedAt = createdAt
	}

	expiresAt := time.Unix(0, redisView.ExpiresAtNs)

	return model.Session{
		UUID:      uuidVal,
		UserUUID:  userUUID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}, nil
}
