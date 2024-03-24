package clothes

import (
	"slices"

	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
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
	err := c.repo.Create(clothes)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClothesUsecase) Update(clothes *domain.Clothes) error {
	return c.repo.Update(clothes)
}

func (c *ClothesUsecase) Get(id utils.UUID) (*domain.Clothes, error) {
	clothes, err := c.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if slices.Equal(clothes.Tags, []string{""}) {
		clothes.Tags = []string{}
	}

	return clothes, nil
}

func (c *ClothesUsecase) Delete(id utils.UUID) error {
	return c.repo.Delete(id)
}

func (c *ClothesUsecase) GetByUser(userID utils.UUID, limit int) ([]domain.Clothes, error) {
	clothes, err := c.repo.GetByUser(userID, limit)
	if err != nil {
		return nil, err
	}

	for i := range clothes {
		if slices.Equal(clothes[i].Tags, []string{""}) {
			clothes[i].Tags = []string{}
		}
	}

	return clothes, nil
}
