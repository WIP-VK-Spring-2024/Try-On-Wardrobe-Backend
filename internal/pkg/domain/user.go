package domain

import "database/sql"

type User struct {
	Model
	Name     string         `gorm:"type:varchar(256)"`
	Email    sql.NullString `gorm:"type:varchar(512)"`
	Password []byte         `gorm:"type:varchar(256)"`
	Gender   Gender         `gorm:"type:gender;default:gender('unknown')"`
}

type UserUsecase interface {
	Create(Credentials) (*User, error)
	GetByName(name string) (*User, error)
}

type UserRepository interface {
	Create(*User) error
	GetByName(name string) (*User, error)
}
