package outfits

import (
	"time"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/outfits"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OutfitUsecase struct {
	repo domain.OutfitRepository
}

func New(repo domain.OutfitRepository) domain.OutfitUsecase {
	return &OutfitUsecase{
		repo: repo,
	}
}

func NewWithSqlcRepo(db *pgxpool.Pool) domain.OutfitUsecase {
	return &OutfitUsecase{
		repo: outfits.New(db),
	}
}

func (u OutfitUsecase) Create(outfit *domain.Outfit) error {
	return u.repo.Create(outfit)
}

func (u OutfitUsecase) Update(outfit *domain.Outfit) error {
	old, err := u.repo.GetById(outfit.ID)
	if err != nil {
		return err
	}

	if outfit.UserID != old.UserID {
		return app_errors.ErrNotOwner
	}

	return u.repo.Update(outfit)
}

func (u OutfitUsecase) Delete(userId, outfitID utils.UUID) error {
	outfit, err := u.repo.GetById(outfitID)
	if err != nil {
		return err
	}

	if outfit.UserID != userId {
		return app_errors.ErrNotOwner
	}

	return u.repo.Delete(outfitID)
}

func (u OutfitUsecase) GetById(id utils.UUID) (*domain.Outfit, error) {
	return u.repo.GetById(id)
}

func (u OutfitUsecase) Get(since time.Time, limit int) ([]domain.Outfit, error) {
	return u.repo.Get(since, limit)
}

func (u OutfitUsecase) GetClothesInfo(outfitId utils.UUID) (map[utils.UUID]string, error) {
	return u.repo.GetClothesInfo(outfitId)
}

func (u OutfitUsecase) GetByUser(id utils.UUID) ([]domain.Outfit, error) {
	return u.repo.GetByUser(id)
}
