package domain

import (
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"
)

//easyjson:json
type Outfit struct {
	Model

	UserID  utils.UUID
	StyleID utils.UUID

	Name       string
	Note       optional.String
	Image      string
	Transforms TransformMap
	Seasons    []Season
}

//easyjson:json
type TransformMap map[utils.UUID]Transform

//easyjson:json
type Transform struct {
	Pos      Vector
	Rotation int
	Scale    float32
}

//easyjson:json
type Vector struct {
	X int
	Y int
}

type OutfitRepository interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(utils.UUID) error
	Get(utils.UUID) (*Outfit, error)
	GetByUser(utils.UUID) ([]Outfit, error)
}
