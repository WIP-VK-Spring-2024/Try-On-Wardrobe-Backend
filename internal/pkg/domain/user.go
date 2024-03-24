package domain

import "try-on/internal/pkg/utils/optional"

type User struct {
	Model
	Name     string          `gorm:"type:varchar(256)"`
	Email    optional.String `gorm:"type:varchar(512)"`
	Password []byte          `gorm:"type:varchar(256)"`
	Gender   Gender          `gorm:"type:gender;default:gender('unknown')"`
}

type UserUsecase interface {
	Create(Credentials) (*User, error)
	GetByName(name string) (*User, error)
}

type UserRepository interface {
	Create(*User) error
	GetByName(name string) (*User, error)
}
