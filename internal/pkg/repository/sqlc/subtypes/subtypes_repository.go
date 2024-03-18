package subtypes

import "try-on/internal/pkg/domain"

type SubtypeRepository struct{}

func (repo *SubtypeRepository) GetAll() ([]domain.Subtype, error) {
	return nil, nil
}
