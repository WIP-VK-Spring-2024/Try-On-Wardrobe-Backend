package domain

import (
	"context"

	"github.com/google/uuid"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingOpts) error
	TryOn(ctx context.Context, opts TryOnOpts) error
	GetTryOnResults() (chan interface{}, error)
}

//easyjson:json
type ClothesProcessingOpts struct {
	UserID    uuid.UUID
	ImageID   uuid.UUID
	FileName  string
	ImageType string
}

//easyjson:json
type TryOnOpts struct {
	UserID          uuid.UUID
	PersonFileName  string
	PersonFilePath  string
	ClothesFileName string
	ClothesFilePath string
}

const (
	ImageTypeCloth    = "cloth"
	ImageTypeFullBody = "full-body"
)
