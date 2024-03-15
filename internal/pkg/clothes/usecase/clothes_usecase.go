package usecase

import (
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"

	"github.com/google/uuid"
)

type ClothesUsecase struct {
	repo domain.ClothesRepository
}

func New(repo domain.ClothesRepository) domain.ClothesUsecase {
	return &ClothesUsecase{
		repo: repo,
	}
}

func (c *ClothesUsecase) Create(clothes *domain.Clothes) error {
	return c.repo.Create(&domain.ClothesModel{
		UserID: clothes.UserID,
		Name:   clothes.Name,
		Type: domain.Type{
			Name: clothes.Type,
		},
	})
}

func (c *ClothesUsecase) Update(clothes *domain.Clothes) error {
	return app_errors.ErrUnimplemented
}

func (c *ClothesUsecase) Get(id uuid.UUID) (*domain.Clothes, error) {
	return nil, app_errors.ErrUnimplemented
}

func (c *ClothesUsecase) Delete(id uuid.UUID) error {
	return app_errors.ErrUnimplemented
}

func (c *ClothesUsecase) GetByUser(userId uuid.UUID, filters *domain.ClothesFilters) ([]domain.Clothes, error) {
	return nil, app_errors.ErrUnimplemented
}
