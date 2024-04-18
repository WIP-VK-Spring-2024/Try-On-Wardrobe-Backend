package domain

import (
	"try-on/internal/pkg/utils"
)

//easyjson:json
type User struct {
	Model

	Name     string `validate:"alphanumunicode | oneof=- _ , . ^ : ; $ # ! + = < > ?"`
	Email    string `validate:"omitempty,email"`
	Password []byte

	Avatar  string
	Gender  Gender  `validate:"omitempty,oneof=male female"`
	Privacy Privacy `validate:"omitempty,oneof=private public friends"`
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
	GetByEmail(email string) (*User, error)
	SearchUsers(name string) ([]User, error)
	GetSubscriptions(utils.UUID) ([]User, error)
}
