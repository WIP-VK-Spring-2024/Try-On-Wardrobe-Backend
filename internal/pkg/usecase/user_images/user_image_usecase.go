package user_images

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/user_images"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserImageUsecase struct {
	repo domain.UserImageRepository
}

func New(db *pgxpool.Pool) domain.UserImageUsecase {
	return &UserImageUsecase{
		repo: user_images.New(db),
	}
}

func (u UserImageUsecase) Create(userImage *domain.UserImage) error {
	err := u.repo.Create(userImage)
	if err != nil {
		return err
	}

	userImage.Image = userImage.Image + "/" + userImage.ID.String()

	return u.repo.SetUserImageUrl(userImage.ID, userImage.Image)
}

func (u UserImageUsecase) SetUserImageUrl(id utils.UUID, url string) error {
	return u.repo.SetUserImageUrl(id, url)
}

func (u UserImageUsecase) Delete(id utils.UUID) error {
	return u.repo.Delete(id)
}

func (u UserImageUsecase) GetByUser(userID utils.UUID) ([]domain.UserImage, error) {
	return u.repo.GetByUser(userID)
}

func (u UserImageUsecase) Get(id utils.UUID) (*domain.UserImage, error) {
	return u.repo.Get(id)
}
