package domain

import "github.com/google/uuid"

type Session struct {
	ID     string
	UserID uuid.UUID
}

//easyjson:json
type Credentials struct {
	Name     string
	Password string
}

type SessionUsecase interface {
	Login(Credentials) (*Session, error)
	Register(user *User) (*Session, error)
	Logout(sessionID string) error
}

type SessionRepository interface {
	Get(key string) (*Session, error)
	Put(session Session) error
	Delete(key string) error
}
