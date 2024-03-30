package domain

//easyjson:json
type Type struct {
	Model
	Name      string
	Tryonable bool
	Subtypes  []Subtype
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
