package users

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.UserRepository {
	return UserRepository{
		db: db,
	}
}

func (repo UserRepository) Create(user *domain.User) error {
	err := repo.db.Create(user).Error
	return utils.GormError(err)
}

func (repo UserRepository) GetByName(name string) (*domain.User, error) {
	user := &domain.User{}
	err := repo.db.First(user, "name = ?", name).Error
	return utils.TranslateGormError(user, err)
}

func (repo UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	err := repo.db.First(user, id).Error
	return utils.TranslateGormError(user, err)
}
