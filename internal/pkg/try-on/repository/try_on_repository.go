package repository

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TryOnResultRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.TryOnResultRepository {
	return &TryOnResultRepository{
		db: db,
	}
}

func (repo *TryOnResultRepository) Create(res *domain.TryOnResult) error {
	err := repo.db.Create(res).Error
	return utils.GormError(err)
}

func (repo *TryOnResultRepository) Delete(id uuid.UUID) error {
	err := repo.db.Delete(&domain.TryOnResult{}, id).Error
	return utils.GormError(err)
}

const initTryOnResNum = 3

func (repo *TryOnResultRepository) GetByUser(userID uuid.UUID) ([]domain.TryOnResult, error) {
	results := make([]domain.TryOnResult, 0, initTryOnResNum)
	queryResult := repo.db.Find(&results, "user_id = ?", userID)

	err := utils.GormError(queryResult.Error)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *TryOnResultRepository) GetByUserAndClothes(userID uuid.UUID, clothesID uuid.UUID) ([]domain.TryOnResult, error) {
	results := make([]domain.TryOnResult, 0, initTryOnResNum)
	queryResult := c.db.Find(&results, "user_id = ?", userID, "clothes_id = ?", clothesID)

	err := utils.GormError(queryResult.Error)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *TryOnResultRepository) Get(ID uuid.UUID) (*domain.TryOnResult, error) {
	res := &domain.TryOnResult{}
	err := c.db.First(res, ID).Error
	return utils.TranslateGormError(res, err)
}

func (repo *TryOnResultRepository) Rate(id uuid.UUID, rating int) error {
	err := repo.db.Model(&domain.TryOnResult{}).Update("rating", rating).Error
	return utils.GormError(err)
}
