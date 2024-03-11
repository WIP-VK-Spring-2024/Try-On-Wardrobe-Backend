package domain

type Tag struct {
	Model
	Name string `gorm:"type:varchar(64);unique"`
}
