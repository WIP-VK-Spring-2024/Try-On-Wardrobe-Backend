package domain

type Style struct {
	Model
	Name string `gorm:"type:varchar(64);unique"`
}
