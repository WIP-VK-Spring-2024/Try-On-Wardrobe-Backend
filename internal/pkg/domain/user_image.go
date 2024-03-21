package domain

import "try-on/internal/pkg/utils"

//easyjson:json
type UserImage struct {
	Model

	UserID utils.UUID
	Image  string
}

type UserImageRepository interface {
	Create(img *UserImage) error
	Get(id utils.UUID) (*UserImage, error)
	Delete(id utils.UUID) error
	GetByUser(userId utils.UUID) ([]UserImage, error)
}
