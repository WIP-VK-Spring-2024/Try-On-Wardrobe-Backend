package outfits

import (
	"time"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type OutfitUsecase struct {
	repo domain.OutfitRepository
}

func New(repo domain.OutfitRepository) domain.OutfitUsecase {
	return &OutfitUsecase{
		repo: repo,
	}
}

func (u *OutfitUsecase) Create(outfit *domain.Outfit) error {
	return u.repo.Create(outfit)
}

func (u *OutfitUsecase) Update(outfit *domain.Outfit) error {
	old, err := u.repo.GetById(outfit.ID)
	if err != nil {
		return err
	}

	if outfit.UserID != old.UserID {
		return app_errors.ErrNotOwner
	}

	return u.repo.Update(outfit)
}

func (u *OutfitUsecase) Delete(userId, outfitID utils.UUID) error {
	outfit, err := u.repo.GetById(outfitID)
	if err != nil {
		return err
	}

	if outfit.UserID != userId {
		return app_errors.ErrNotOwner
	}

	return u.repo.Delete(outfitID)
}

func (u *OutfitUsecase) GetById(id utils.UUID) (*domain.Outfit, error) {
	return u.repo.GetById(id)
}

func (u *OutfitUsecase) Get(since time.Time, limit int) ([]domain.Outfit, error) {
	return u.repo.Get(since, limit)
}

func (u *OutfitUsecase) GetByUser(id utils.UUID) ([]domain.Outfit, error) {
	return u.repo.GetByUser(id)
}
