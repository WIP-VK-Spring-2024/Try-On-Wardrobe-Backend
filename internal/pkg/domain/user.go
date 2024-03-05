package domain

import "try-on/internal/pkg/domain/gender"

type User struct {
	Model
	Name     string
	Email    string
	Password string
	Gender   gender.Gender
}

type UserRepository interface {
	Create(*User) error
}
