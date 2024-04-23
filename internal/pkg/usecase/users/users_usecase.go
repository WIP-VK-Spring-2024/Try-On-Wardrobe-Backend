package users

import (
	"log"

	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type UserUsecase struct {
	repo   domain.UserRepository
	images domain.UserImageRepository
	cfg    config.DefaultImgPaths
}

func New(repo domain.UserRepository, images domain.UserImageRepository, cfg config.DefaultImgPaths) domain.UserUsecase {
	return &UserUsecase{
		repo:   repo,
		images: images,
		cfg:    cfg,
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

	if user.Gender != domain.Male && user.Gender != domain.Female {
		user.Gender = domain.Female
	}

	err = u.images.Create(&domain.UserImage{
		UserID: user.ID,
		Image:  u.cfg[user.Gender],
	})
	if err != nil {
		log.Println("Error creating default image:", err)
	}

	return nil
}

func (u UserUsecase) Update(user domain.User) error {
	return u.repo.Update(user)
}

func (u UserUsecase) GetByName(name string) (*domain.User, error) {
	return u.repo.GetByName(name)
}

func (u UserUsecase) SearchUsers(opts domain.SearchUserOpts) ([]domain.User, error) {
	if opts.Limit == 0 {
		opts.Limit = 16
	}

	return u.repo.SearchUsers(opts)
}

func (u UserUsecase) GetSubscriptions(id utils.UUID) ([]domain.User, error) {
	return u.repo.GetSubscriptions(id)
}
