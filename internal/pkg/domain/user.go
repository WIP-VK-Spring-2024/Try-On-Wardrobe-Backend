package domain

type User struct {
	Model
	Name     string `gorm:"type:varchar(256)"`
	Email    string `gorm:"type:varchar(512);unique"`
	Password []byte `gorm:"type:varchar(256)"`
	Gender   Gender `gorm:"type:gender"`
}

type UserRepository interface {
	Create(*User) error
	GetByName(name string) (*User, error)
}
