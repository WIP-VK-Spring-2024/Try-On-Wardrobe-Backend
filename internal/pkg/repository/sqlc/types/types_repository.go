package types

import "try-on/internal/pkg/domain"

type TypeRepository struct{}

func (repo *TypeRepository) GetAll() ([]domain.Type, error) {
	return nil, nil
}
