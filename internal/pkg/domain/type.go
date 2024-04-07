package domain

//easyjson:json
type Type struct {
	Model
	Name      string
	Tryonable bool `json:"tryonable,!omitempty"` //lint:ignore SA5008 easyjson custom tags
	Subtypes  []Subtype
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
