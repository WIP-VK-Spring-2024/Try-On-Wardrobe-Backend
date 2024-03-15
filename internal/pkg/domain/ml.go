package domain

import (
	"context"

	"github.com/google/uuid"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingOpts) error
}

//easyjson:json
type ClothesProcessingOpts struct {
	UserID    uuid.UUID
	ImageID   uuid.UUID
	FileName  string
	ImageType string
}

const (
	ImageTypeCloth    = "cloth"
	ImageTypeFullBody = "full-body"
)
