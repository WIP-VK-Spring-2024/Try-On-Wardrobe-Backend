package domain

import (
	"context"

	"try-on/internal/pkg/utils"

	"go.uber.org/zap"
)

const (
	TryOnCategoryDress = "dresses"
	TryOnCategoryUpper = "upper_body"
	TryOnCategoryLower = "lower_body"
)

type ClothesProcessingModel interface {
	Process(ctx context.Context, opts ClothesProcessingRequest) error
	TryOn(ctx context.Context, opts TryOnRequest) error
	GetTryOnResults(logger *zap.SugaredLogger, handler func(*TryOnResponse) Result) error
	GetProcessingResults(logger *zap.SugaredLogger, handler func(*ClothesProcessingResponse) Result) error
	Close()
}

type ClothesClassificationRepository interface {
	GetClassifications(userId utils.UUID, tagLimit int32) (*ClothesClassificationRequest, error)
	GetTypeId(engName string) (utils.UUID, error)
	GetSubtypeIds(engNames []string) ([]utils.UUID, error)
	GetStyleId(engName string) (utils.UUID, error)
	GetTags(engNames []string) ([]string, error)
}

//easyjson:json
type ClothesProcessingRequest struct {
	UserID         utils.UUID
	ClothesID      utils.UUID
	ClothesDir     string
	Classification ClothesClassificationRequest
}

//easyjson:json
type ClothesProcessingResponse struct {
	UserID         utils.UUID
	ClothesID      utils.UUID
	ResultDir      string
	Classification ClothesClassificationResponse
}

//easyjson:json
type ClothesClassificationRequest struct { // Request to ML-server
	Tags          []string
	Styles        []string
	Categories    []string
	Subcategories []string
	Seasons       []string
}

//easyjson:json
type ClothesClassificationResponse struct { // End-user response
	Type     utils.UUID
	Subtypes []utils.UUID // maybe only one should be returned?
	Style    utils.UUID
	Seasons  []Season
	Tags     []string
}

//easyjson:json
type TryOnRequest struct {
	UserID       utils.UUID
	UserImageID  utils.UUID
	UserImageDir string
	ClothesDir   string
	Clothes      []TryOnClothesInfo
}

//easyjson:json
type TryOnClothesInfo struct {
	ClothesID utils.UUID
	Category  string
}

//easyjson:json
type TryOnResponse struct {
	UserID         utils.UUID
	ClothesID      []utils.UUID
	UserImageID    utils.UUID
	TryOnResultID  string
	TryOnResultDir string
}
