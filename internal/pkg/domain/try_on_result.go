package domain

import "try-on/internal/pkg/utils"

type TryOnResultRepository interface {
	Create(res *TryOnResult) error
	Delete(id utils.UUID) error
	GetByUser(userID utils.UUID) ([]TryOnResult, error)
	GetLast(userID utils.UUID) (*TryOnResult, error)
	GetByClothes(clothesID utils.UUID) ([]TryOnResult, error)
	Rate(id utils.UUID, rating int) error
}

//easyjson:json
type TryOnResult struct {
	Model

	Image       string
	Rating      int
	UserImageID utils.UUID
	ClothesID   utils.UUID
}
