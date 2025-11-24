package model

import (
	"time"

	"github.com/google/uuid"
)

// User представляет доменную модель пользователя.
type User struct {
	UUID      uuid.UUID
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo содержит базовую информацию о пользователе.
type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

// NotificationMethod определяет способ уведомления пользователя.
type NotificationMethod struct {
	ProviderName string // Провайдер: telegram, email, push и т.д.
	Target       string // Адрес/идентификатор назначения (email, чат-id)
}
