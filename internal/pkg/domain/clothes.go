package domain

import (
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"
)

//easyjson:json
type Clothes struct {
	Model

	Name string
	Note optional.String
	Tags []string

	UserID utils.UUID
	Image  string

	StyleID utils.UUID
	Style   string `json:"-"`

	TypeID utils.UUID
	Type   string `json:"-"`

	SubtypeID utils.UUID
	Subtype   string `json:"-"`

	Color   optional.String
	Seasons []Season
}

type ClothesUsecase interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	SetImage(id utils.UUID, path string) error
	Get(id utils.UUID) (*Clothes, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]Clothes, error)
}

type ClothesRepository interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	SetImage(id utils.UUID, path string) error
	Get(id utils.UUID) (*Clothes, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]Clothes, error)
}
