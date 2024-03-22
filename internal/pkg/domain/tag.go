package domain

//easyjson:json
type Tag struct {
	Model
	Name     string
	UseCount int32
}

type TagRepository interface {
	Get(limit, offset int) ([]Tag, error)
}
