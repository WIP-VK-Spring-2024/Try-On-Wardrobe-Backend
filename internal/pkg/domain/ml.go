package domain

import (
	"context"

	"try-on/internal/pkg/utils"

	"go.uber.org/zap"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingRequest) error
	TryOn(ctx context.Context, opts TryOnRequest) error
	GetTryOnResults(logger *zap.SugaredLogger, handler func(*TryOnResponse) Result) error
	GetProcessingResults(logger *zap.SugaredLogger, handler func(*ClothesProcessingResponse) Result) error
	Close()
}

//easyjson:json
type ClothesProcessingRequest struct {
	UserID     utils.UUID
	ClothesID  utils.UUID
	ClothesDir string
}

//easyjson:json
type ClothesProcessingResponse struct {
	UserID       utils.UUID
	ClothesID    utils.UUID
	ProcessedDir string
}

//easyjson:json
type TryOnRequest struct {
	UserID       utils.UUID
	UserImageID  utils.UUID
	ClothesID    utils.UUID
	UserImageDir string
	ClothesDir   string
}

//easyjson:json
type TryOnResponse struct {
	UserID         utils.UUID
	ClothesID      utils.UUID
	UserImageID    utils.UUID
	TryOnResultID  string
	TryOnResultDir string
}
