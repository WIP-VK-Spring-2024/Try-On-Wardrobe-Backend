package session

import "github.com/google/uuid"

type Session struct {
	Key    string
	UserID uuid.UUID
}

type SessionRepository interface {
	Get(key string) (uuid.UUID, error)
	Put(session *Session) error
	Delete(key string) error
}
