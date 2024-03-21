package domain

//easyjson:json
type Style struct {
	Model
	Name string `gorm:"type:varchar(64)"`
}

type StylesRepository interface {
	GetAll() ([]Style, error)
}
