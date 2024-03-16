package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type ClothesModel struct {
	Model

	Name string         `gorm:"type:varchar(128)"`
	Note sql.NullString `gorm:"type:varchar(512)"`
	Tags []Tag          `gorm:"many2many:clothes_tags;"`

	Image string `gorm:"type:varchar(256)"`

	UserID uuid.UUID
	User   User

	StyleID uuid.UUID `gorm:"default:null"`
	Style   *Style

	TypeID uuid.UUID `gorm:"default:null"`
	Type   Type

	SubtypeID uuid.UUID `gorm:"default:null"`
	Subtype   Subtype

	Color   sql.NullString `gorm:"type:char(7)"`
	Seasons []Season       `gorm:"type:season[]"`
}

type Clothes struct {
	ID   uuid.UUID
	Name string
	Note string
	Tags []string

	Image string

	UserID uuid.UUID

	Style   string
	Type    string
	Subtype string

	Color   string
	Seasons []Season
}

func (ClothesModel) TableName() string {
	return "clothes"
}

type ClothesFilters struct {
	Limit   int
	Tags    []string
	Style   string
	Type    string
	Subtype string
	Color   uint32
	Seasons []Season
}

type ClothesUsecase interface {
	Create(clothes *Clothes) error
	Update(clothes *Clothes) error
	Get(id uuid.UUID) (*Clothes, error)
	Delete(id uuid.UUID) error
	GetByUser(userId uuid.UUID, filters *ClothesFilters) ([]Clothes, error)
}

type ClothesRepository interface {
	Create(clothes *ClothesModel) error
	Update(clothes *ClothesModel) error
	Get(id uuid.UUID) (*ClothesModel, error)
	Delete(id uuid.UUID) error
	GetByUser(userId uuid.UUID, filters *ClothesFilters) ([]ClothesModel, error)
}
