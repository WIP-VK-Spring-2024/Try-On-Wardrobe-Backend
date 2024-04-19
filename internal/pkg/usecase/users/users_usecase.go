package users

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type UserUsecase struct {
	repo domain.UserRepository
}

func New(repo domain.UserRepository) domain.UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (u UserUsecase) Create(user *domain.User) error {
	salt, err := utils.NewSalt()
	if err != nil {
		return err
	}

	user.Password = string(utils.Hash([]byte(user.Password), salt)) + ":" + string(salt)

	err = u.repo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (u UserUsecase) Update(user domain.User) error {
	return u.repo.Update(user)
}

func (u UserUsecase) GetByName(name string) (*domain.User, error) {
	return u.repo.GetByName(name)
}

func (u UserUsecase) SearchUsers(name string) ([]domain.User, error) {
	return u.repo.SearchUsers(name)
}

func (u UserUsecase) GetSubscriptions(id utils.UUID) ([]domain.User, error) {
	return u.repo.GetSubscriptions(id)
}
