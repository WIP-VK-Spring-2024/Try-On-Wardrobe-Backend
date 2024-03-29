package domain

//easyjson:json
type Style struct {
	Model
	Name string
}

type StyleRepository interface {
	GetAll() ([]Style, error)
}
