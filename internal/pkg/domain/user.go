package domain

import (
	"try-on/internal/pkg/utils"
)

//easyjson:json
type User struct {
	Model

	Name     string `validate:"required,alphanum"`
	Email    string `validate:"required,email"`
	Password string

	Avatar  string
	Gender  Gender  `validate:"omitempty,oneof=male female"`
	Privacy Privacy `validate:"omitempty,oneof=private public friends"`
}

type SearchUserOpts struct {
	UserID utils.UUID
	Name   string `query:"name"`
	Limit  int    `query:"limit"`
	Since  string `query:"since"`
}

type UserUsecase interface {
	Create(user *User) error
	Update(user User) error
	GetByName(name string) (*User, error)
	SearchUsers(opts SearchUserOpts) ([]User, error)
	GetSubscriptions(utils.UUID) ([]User, error)
}

type UserRepository interface {
	Create(*User) error
	Update(user User) error
	GetByName(name string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id utils.UUID) (*User, error)
	SearchUsers(opts SearchUserOpts) ([]User, error)
	GetSubscriptions(utils.UUID) ([]User, error)
}
