package domain

import (
	"context"

	"github.com/google/uuid"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingOpts) error
	TryOn(ctx context.Context, opts TryOnOpts) error
	GetTryOnResults() (chan TryOnResponse, error)
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
	ClothesID       uuid.UUID
	PersonFileName  string
	PersonFilePath  string
	ClothesFileName string
	ClothesFilePath string
}

//easyjson:json
type TryOnResponse struct {
	UserID      uuid.UUID
	ClothesID   uuid.UUID
	ResFileName string
	ResFilePath string
}

const (
	ImageTypeCloth    = "cloth"
	ImageTypeFullBody = "full-body"
)
