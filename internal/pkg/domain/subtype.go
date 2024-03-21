package domain

import "try-on/internal/pkg/utils"

//easyjson:json
type Subtype struct {
	Model
	Name   string `gorm:"type:varchar(64)"`
	TypeID utils.UUID
}

type SubtypeRepository interface {
	GetAll() ([]Subtype, error)
}
