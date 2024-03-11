package domain

type Subtype struct {
	Model
	Name string `gorm:"type:varchar(64);unique"`
}
