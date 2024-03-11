package usecase

import (
	"slices"

	"try-on/internal/pkg/app_errors"
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

func (u *UserUsecase) Create(creds domain.Credentials) (*domain.User, error) {
	salt, err := utils.NewSalt()
	if err != nil {
		return nil, app_errors.New(err)
	}

	user := &domain.User{
		Name:     creds.Name,
		Password: slices.Concat(utils.Hash([]byte(creds.Password), salt), []byte{':'}, salt),
	}

	err = u.repo.Create(user)
	if err != nil {
		return nil, app_errors.New(err)
	}

	return user, nil
}

func (u *UserUsecase) GetByName(name string) (*domain.User, error) {
	return nil, nil
}
