package model

import (
	"time"

	"github.com/google/uuid"
)

// Session представляет доменную модель пользовательской сессии.
type Session struct {
	UUID      uuid.UUID
	UserUUID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
