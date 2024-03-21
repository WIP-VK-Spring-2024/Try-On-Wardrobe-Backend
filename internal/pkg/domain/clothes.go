package domain

import (
	"database/sql"

	"try-on/internal/pkg/utils"
)

type ClothesModel struct {
	Model

	Name string         `gorm:"type:varchar(128)"`
	Note sql.NullString `gorm:"type:varchar(512)"`
	Tags []Tag          `gorm:"many2many:clothes_tags;"`

	Image string `gorm:"type:varchar(256)"`

	UserID utils.UUID
	User   User

	StyleID utils.UUID `gorm:"default:null"`
	Style   *Style

	TypeID utils.UUID `gorm:"default:null"`
	Type   Type

	SubtypeID utils.UUID `gorm:"default:null"`
	Subtype   Subtype

	Color   sql.NullString `gorm:"type:char(7)"`
	Seasons []Season       `gorm:"type:season[]"`
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
