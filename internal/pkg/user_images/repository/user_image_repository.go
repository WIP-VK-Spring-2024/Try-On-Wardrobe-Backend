package repository

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserImageRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.UserImageRepository {
	return &UserImageRepository{
		db: db,
	}
}

func (repo *UserImageRepository) Create(userImage *domain.UserImage) error {
	err := repo.db.Create(userImage).Error
	return utils.GormError(err)
}

func (repo *UserImageRepository) Delete(id uuid.UUID) error {
	err := repo.db.Delete(&domain.UserImage{}, id).Error
	return utils.GormError(err)
}

const initUserImageNum = 3

func (repo *UserImageRepository) GetByUser(userID uuid.UUID) ([]domain.UserImage, error) {
	userImages := make([]domain.UserImage, 0, initUserImageNum)
	queryResult := repo.db.Find(&userImages, "user_id = ?", userID)

	err := utils.GormError(queryResult.Error)
	if err != nil {
		return nil, err
	}
	return userImages, nil
}

func (c *UserImageRepository) Get(id uuid.UUID) (*domain.UserImage, error) {
	userImage := &domain.UserImage{}
	err := c.db.First(userImage, id).Error
	return utils.TranslateGormError(userImage, err)
}
