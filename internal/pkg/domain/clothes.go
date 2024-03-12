package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type Clothes struct {
	Model

	Name string         `gorm:"type:varchar(128)"`
	Note sql.NullString `gorm:"type:varchar(512)"`
	Tags []Tag          `gorm:"many2many:clothes_tags;"`

	StyleID uuid.UUID
	Style   *Style

	TypeID uuid.UUID
	Type   Type

	SubtypeID uuid.UUID
	Subtype   Subtype

	Color   sql.NullString `gorm:"type:char(7)"`
	Seasons []Season       `gorm:"type:season[]"`
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
