package domain

import (
	"try-on/internal/pkg/utils"
)

//easyjson:json
type User struct {
	Model

	Name     string
	Email    string
	Password []byte

	Avatar  string
	Gender  Gender
	Privacy Privacy
}

type UserUsecase interface {
	Create(user *User) error
	Update(user User) error
	GetByName(name string) (*User, error)
	SearchUsers(name string) ([]User, error)
	GetSubscriptions(utils.UUID) ([]User, error)
}

type UserRepository interface {
	Create(*User) error
	Update(user User) error
	GetByName(name string) (*User, error)
	SearchUsers(name string) ([]User, error)
	GetSubscriptions(utils.UUID) ([]User, error)
}
