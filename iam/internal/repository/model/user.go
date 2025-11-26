package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID                uuid.UUID `db:"uuid"`
	Login               string    `db:"login"`
	Email               string    `db:"email"`
	PasswordHash        string    `db:"password_hash"`
	NotificationMethods []byte    `db:"notification_methods"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}
