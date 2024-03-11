package domain

import "database/sql"

type User struct {
	Model
	Name     string         `gorm:"type:varchar(256);uniqueIndex:uni_users_name,class:varchar_pattern_ops"`
	Email    sql.NullString `gorm:"type:varchar(512);unique"`
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
