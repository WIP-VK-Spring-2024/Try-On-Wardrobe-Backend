package domain

import "try-on/internal/pkg/utils"

//easyjson:json
type UserImage struct {
	Model

	UserID utils.UUID
	Image  string
}

type UserImageUsecase interface {
	Create(img *UserImage) error
	Get(id utils.UUID) (*UserImage, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID) ([]UserImage, error)
}

type UserImageRepository interface {
	Create(img *UserImage) error
	SetUserImageUrl(id utils.UUID, url string) error
	Get(id utils.UUID) (*UserImage, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID) ([]UserImage, error)
}
