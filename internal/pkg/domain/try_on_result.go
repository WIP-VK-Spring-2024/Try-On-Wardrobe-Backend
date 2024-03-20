package domain

import "github.com/google/uuid"

type TryOnResultRepository interface {
	Create(res *TryOnResult) error
	Delete(id uuid.UUID) error
	GetByUser(userID uuid.UUID) ([]TryOnResult, error)
	GetLast(userID uuid.UUID) (*TryOnResult, error)
	GetByClothes(clothesID uuid.UUID) ([]TryOnResult, error)
	Rate(id uuid.UUID, rating int) error
}

type TryOnResult struct {
	Model

	Image       string
	Rating      int
	UserImageID uuid.UUID
	ClothesID   uuid.UUID
}
