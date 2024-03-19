package domain

//easyjson:json
type Type struct {
	Model
	Name string `gorm:"type:varchar(64)"`
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
