package domain

import (
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"
)

const (
	ClothesStatusCreated   = "created"
	ClothesStatusProcessed = "processed"
)

//easyjson:json
type Clothes struct {
	Model

	Name      string `validate:"name"`
	Tryonable bool   `json:"tryonable,!omitempty"` //lint:ignore SA5008 easyjson custom tags
	Note      optional.String
	Tags      []string `validate:"dive,name"`

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

	Privacy Privacy
}

type ClothesUsecase interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	SetImage(id utils.UUID, path string) error
	Get(id utils.UUID) (*Clothes, error)
	GetTryOnInfo(ids []utils.UUID) ([]TryOnClothesInfo, error)
	Delete(userId, clothesId utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]Clothes, error)
}

type ClothesRepository interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	SetImage(id utils.UUID, path string) error
	Get(id utils.UUID) (*Clothes, error)
	GetTryOnInfo(ids []utils.UUID) ([]TryOnClothesInfo, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]Clothes, error)
	GetByWeather(userId utils.UUID, temp *int) ([]GenClothesInfo, error)
}
