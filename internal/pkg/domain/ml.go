package domain

import "github.com/google/uuid"

type ClothesProcessingModel interface {
	Process(opts ClothesProcessingOpts) error
}

type ClothesProcessingOpts struct {
	UserID        uuid.UUID
	ClothesID     uuid.UUID
	ImagePath     string
	CutBackground bool
	Categorise    bool
}
