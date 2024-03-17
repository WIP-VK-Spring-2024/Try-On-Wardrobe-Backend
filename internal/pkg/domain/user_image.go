package domain

import "github.com/google/uuid"

type UserImage struct {
	Model

	UserID uuid.UUID
	Image  string
}

type UserImageRepository interface {
	Create(img *UserImage) error
	Get(id uuid.UUID) (*UserImage, error)
	Delete(id uuid.UUID) error
	GetByUser(userId uuid.UUID) ([]UserImage, error)
}
