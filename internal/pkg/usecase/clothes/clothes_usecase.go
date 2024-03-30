package clothes

import (
	"slices"

	"try-on/internal/pkg/app_errors"
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
	return c.repo.Create(clothes)
}

func (c *ClothesUsecase) Update(clothes *domain.Clothes) error {
	old, err := c.repo.Get(clothes.ID)
	if err != nil {
		return err
	}

	if old.UserID != clothes.UserID {
		return app_errors.ErrNotOwner
	}

	return c.repo.Update(clothes)
}

func (c *ClothesUsecase) SetImage(id utils.UUID, path string) error {
	return c.repo.SetImage(id, path)
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

func (c *ClothesUsecase) Delete(userId, id utils.UUID) error {
	clothes, err := c.repo.Get(id)
	if err != nil {
		return err
	}

	if clothes.UserID != userId {
		return app_errors.ErrNotOwner
	}
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
