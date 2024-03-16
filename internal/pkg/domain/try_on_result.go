package domain

import "github.com/google/uuid"

type TryOnResult struct {
	Model

	Image  string
	Rating int

	UserID uuid.UUID
	USer   *User

	ClothesModelID uuid.UUID
	ClothesModel   *ClothesModel
}
