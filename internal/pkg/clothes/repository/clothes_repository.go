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
	err := c.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Table: "tags", Name: "name"},
			{Table: "styles", Name: "name"},
			{Table: "types", Name: "name"},
			{Table: "subtypes", Name: "name"},
		},
		DoNothing: true,
	}).Create(clothes).Error

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
	err := utils.GormError(c.db.Find(clothes, "user_id = ?", userID).Error)
	if err != nil {
		return nil, err
	}
	return clothes, nil
}
