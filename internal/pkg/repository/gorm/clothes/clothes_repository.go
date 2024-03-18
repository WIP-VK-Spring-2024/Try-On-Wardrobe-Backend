package clothes

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
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"name"}),
			})
		}

		if clothes.Style != nil {
			if err := clauses().Create(clothes.Style).Error; err != nil {
				return err
			}
		}

		// if err := clauses().Create(&clothes.Subtype).Error; err != nil {
		// 	return err
		// }

		// if err := clauses().Create(&clothes.Type).Error; err != nil {
		// 	return err
		// }

		tagsTmp := make([]domain.Tag, 0, len(clothes.Tags))
		copy(tagsTmp, clothes.Tags)
		clothes.Tags = nil

		err := tx.Debug().Create(clothes).Error
		if err != nil {
			return err
		}

		// return nil
		return tx.Debug().Model(clothes).Association("Tags").Append(tagsTmp)
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
	err := c.db.Preload(clause.Associations).First(clothes, id).Error
	return utils.TranslateGormError(clothes, err)
}

func (c *ClothesRepository) Delete(id uuid.UUID) error {
	err := c.db.Delete(&domain.ClothesModel{}, id).Error
	return utils.GormError(err)
}

func (c *ClothesRepository) GetByUser(userID uuid.UUID, limit int) ([]domain.ClothesModel, error) {
	clothes := make([]domain.ClothesModel, 0, initClothesNum)

	query := c.db
	if limit != 0 {
		query = query.Limit(limit)
	}

	result := query.Preload(clause.Associations).Find(&clothes, "user_id = ?", userID)

	err := utils.GormError(result.Error)
	if err != nil {
		return nil, err
	}
	return clothes, nil
}
