package domain

import (
	"github.com/google/uuid"
)

type Clothes struct {
	Model

	Name string `gorm:"type:varchar(128)"`
	Note string `gorm:"type:varchar(512)"`
	Tags []Tag  `gorm:"many2many:clothes_tags;"`

	StyleID int
	Style   Style

	TypeID int
	Type   Type

	SubtypeID int
	Subtype   Subtype

	Color   string   `gorm:"type:char(7)"`
	Seasons []Season `gorm:"type:season[]"`
}

type ClothesFilters struct {
	Tags    []string
	Style   string
	Type    string
	Subtype string
	Color   uint32
	Seasons []Season
}

type ClothesRepository interface {
	Get(id uuid.UUID) (*Clothes, error)
	Delete(id uuid.UUID) error
	GetByUser(userId uuid.UUID, filters *ClothesFilters) ([]Clothes, error)
	Update(clothes *Clothes) error
}
