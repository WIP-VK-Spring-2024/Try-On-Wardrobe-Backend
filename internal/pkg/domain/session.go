package domain

import "try-on/internal/pkg/utils"

type Session struct {
	ID     string
	UserID utils.UUID
}

//easyjson:json
type Credentials struct {
	Name     string
	Password string
}

type SessionUsecase interface {
	Login(Credentials) (*Session, error)
	IsLoggedIn(*Session) (bool, error)
	IssueToken(id utils.UUID) (string, error)
}

type SessionRepository interface {
	Get(key string) (*Session, error)
	Put(session Session) error
	Delete(key string) error
}
