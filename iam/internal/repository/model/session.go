package model

type SessionRedisView struct {
	SessionUUID string `redis:"session_uuid"`
	UserUUID    string `redis:"user_uuid"`
	CreatedAt   int64  `redis:"created_at"`
	UpdatedAt   int64  `redis:"updated_at"`
	ExpiresAt   int64  `redis:"expires_at"`
}
