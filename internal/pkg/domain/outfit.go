package domain

import (
	"time"

	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"
)

//easyjson:json
type Outfit struct {
	Model

	UserID  utils.UUID
	StyleID utils.UUID

	Public bool

	Name       string
	Note       optional.String
	Image      string
	Transforms TransformMap
	Seasons    []Season
	Tags       []string
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

type OutfitUsecase interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(userId, outfitID utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) (map[utils.UUID]string, error)
	GetByUser(utils.UUID) ([]Outfit, error)
}

type OutfitRepository interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) (map[utils.UUID]string, error)
	GetByUser(utils.UUID) ([]Outfit, error)
}
