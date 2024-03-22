package domain

//easyjson:json
type Style struct {
	Model
	Name string
}

type StylesRepository interface {
	GetAll() ([]Style, error)
}
