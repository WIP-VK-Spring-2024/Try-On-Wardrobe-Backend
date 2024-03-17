package domain

import "github.com/google/uuid"

type TryOnResultRepository interface {
	Create(res *TryOnResult) error
	Delete(id uuid.UUID) error
	GetByUser(userID uuid.UUID) ([]TryOnResult, error)
	Get(id uuid.UUID) (*TryOnResult, error)
	GetByUserAndClothes(userID uuid.UUID, clothesID uuid.UUID) ([]TryOnResult, error)
	Rate(id uuid.UUID, rating int) error
}

type TryOnResult struct {
	Model

	Image  string
	Rating int

	UserID uuid.UUID
	User   *User

	ClothesModelID uuid.UUID
	ClothesModel   *ClothesModel
}
