package domain

import "github.com/google/uuid"

//easyjson:json
type Subtype struct {
	Model
	Name string `gorm:"type:varchar(64)"`

	TypeID uuid.UUID
	Type   *Type
}

type SubtypeRepository interface {
	GetAll() ([]Subtype, error)
}
