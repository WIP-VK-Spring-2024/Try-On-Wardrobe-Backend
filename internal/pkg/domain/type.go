package domain

type Type struct {
	Model
	Name string `gorm:"type:varchar(64)"`
}
