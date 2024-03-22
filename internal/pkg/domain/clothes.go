package domain

import (
	"database/sql"

	"try-on/internal/pkg/utils"
)

type ClothesModel struct {
	Model

	Name string
	Note sql.NullString
	Tags []Tag

	Image string

	UserID utils.UUID
	User   User

	StyleID utils.UUID
	Style   *Style

	TypeID utils.UUID
	Type   Type

	SubtypeID utils.UUID
	Subtype   Subtype

	Color   sql.NullString
	Seasons []Season
}

func (*ClothesModel) TableName() string {
	return "clothes"
}

//easyjson:json
type Clothes struct {
	ID   utils.UUID
	Name string
	Note string
	Tags []string

	Image string

	UserID utils.UUID

	Style   string
	Type    string
	Subtype string

	Color   string
	Seasons []Season
}

type ClothesUsecase interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	Get(id utils.UUID) (*Clothes, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]Clothes, error)
}

type ClothesRepository interface {
	Create(clothes *ClothesModel) error
	Update(clothes *ClothesModel) error
	Get(id utils.UUID) (*ClothesModel, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID, limit int) ([]ClothesModel, error)
}
