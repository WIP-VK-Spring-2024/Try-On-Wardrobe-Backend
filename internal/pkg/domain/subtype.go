package domain

import "try-on/internal/pkg/utils"

//easyjson:json
type Subtype struct {
	Model
	Name   string
	TypeID utils.UUID
}

type SubtypeRepository interface {
	GetAll() ([]Subtype, error)
}
