package domain

//easyjson:json
type Type struct {
	Model
	Name     string
	Subtypes []Subtype
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
