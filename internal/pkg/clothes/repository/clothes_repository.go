package repository

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClothesRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.ClothesRepository {
	return &ClothesRepository{
		db: db,
	}
}

const initClothesNum = 15

func (c *ClothesRepository) Create(clothes *domain.ClothesModel) error {
	err := c.db.Transaction(func(tx *gorm.DB) error {
		clauses := func() *gorm.DB {
			return tx.Clauses(clause.OnConflict{DoNothing: true})
		}

		if clothes.Style != nil {
			if err := clauses().Create(clothes.Style).Error; err != nil {
				return err
			}
		}

		if err := clauses().Create(&clothes.Subtype).Error; err != nil {
			return err
		}

		if err := clauses().Create(&clothes.Type).Error; err != nil {
			return err
		}

		return tx.Create(clothes).Error
	},
	)

	return utils.GormError(err)
}

func (c *ClothesRepository) Update(clothes *domain.ClothesModel) error {
	err := c.db.Updates(clothes).Error
	return utils.GormError(err)
}

func (c *ClothesRepository) Get(id uuid.UUID) (*domain.ClothesModel, error) {
	clothes := &domain.ClothesModel{}
	err := c.db.First(clothes, id).Error
	return utils.TranslateGormError(clothes, err)
}

func (c *ClothesRepository) Delete(id uuid.UUID) error {
	err := c.db.Delete(&domain.ClothesModel{}, id).Error
	return utils.GormError(err)
}

func (c *ClothesRepository) GetByUser(userID uuid.UUID, filters *domain.ClothesFilters) ([]domain.ClothesModel, error) {
	clothes := make([]domain.ClothesModel, 0, initClothesNum)

	result := c.db.Limit(initClothesNum).Find(&clothes, "user_id = ?", userID)

	err := utils.GormError(result.Error)
	if err != nil {
		return nil, err
	}
	return clothes, nil
}

func (c *ClothesRepository) GetTryOnResult(userID uuid.UUID, clothesID uuid.UUID) (*domain.TryOnResult, error) {
	res := &domain.TryOnResult{}
	err := c.db.First(res, userID, clothesID).Error
	return utils.TranslateGormError(res, err)
}
