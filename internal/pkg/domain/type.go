package domain

//easyjson:json
type Type struct {
	Model
	Name     string `gorm:"type:varchar(64)"`
	Subtypes []Subtype
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
