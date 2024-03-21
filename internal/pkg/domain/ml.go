package domain

import (
	"context"

	"try-on/internal/pkg/utils"

	"go.uber.org/zap"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingOpts) error
	TryOn(ctx context.Context, opts TryOnOpts) error
	GetTryOnResults(logger *zap.SugaredLogger, handler func(*TryOnResponse) Result) error
	Close()
}

//easyjson:json
type ClothesProcessingOpts struct {
	UserID    utils.UUID
	ImageID   utils.UUID
	FileName  string
	ImageType string
}

type ClothesProcessingResponse struct{}

//easyjson:json
type TryOnOpts struct {
	UserImageID     utils.UUID
	UserID          utils.UUID
	ClothesID       utils.UUID
	PersonFileName  string
	PersonFilePath  string
	ClothesFileName string
	ClothesFilePath string
}

//easyjson:json
type TryOnResponse struct {
	UserID      utils.UUID
	ClothesID   utils.UUID
	UserImageID utils.UUID
	ResFileName string
	ResFilePath string
}

const (
	ImageTypeCloth    = "cloth"
	ImageTypeFullBody = "full-body"
)
