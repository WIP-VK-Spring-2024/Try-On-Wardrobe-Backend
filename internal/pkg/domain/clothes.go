package domain

import (
	seasons "try-on/internal/pkg/domain/season"

	"github.com/google/uuid"
)

type Clothes struct {
	Model
	Name    string
	Note    string
	Tags    []string
	Style   string
	Type    string
	Subtype string
	Color   uint32
	Seasons []seasons.Season
}

type ClothesFilters struct {
	Tags    []string
	Style   string
	Type    string
	Subtype string
	Color   uint32
	Seasons []seasons.Season
}

type ClothesRepository interface {
	Get(id uuid.UUID) (*Clothes, error)
	Delete(id uuid.UUID) error
	GetByUser(userId uuid.UUID, filters *ClothesFilters) ([]Clothes, error)
	Update(clothes *Clothes) error
}
