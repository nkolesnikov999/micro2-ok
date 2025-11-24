package model

type SessionRedisView struct {
	UUID        string `redis:"uuid"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at,omitempty"`
	ExpiresAtNs int64  `redis:"expires_at"`
}
