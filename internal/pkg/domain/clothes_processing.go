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
	Closer
	Process(ctx context.Context, opts ClothesProcessingRequest) error
	GetProcessingResults(logger *zap.SugaredLogger, handler func(*ClothesProcessingResponse) Result) error
}

type ClothesClassificationRepository interface {
	GetTypeBySubtype(subtypeId utils.UUID) (utils.UUID, bool, error)
	GetClassifications(userId utils.UUID, tagLimit int32) (*ClothesClassificationRequest, error)
	GetTypeId(engName string) (utils.UUID, error)
	GetSubtypeIds(engName string) (utils.UUID, error)
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
	QueueResponse

	UserID         utils.UUID
	ClothesID      utils.UUID
	ClothesDir     string
	Tryonable      bool
	Classification ClothesClassificationResponse
}

//easyjson:json
type ClothesProcessingModelResponse struct {
	QueueResponse

	UserID         utils.UUID
	ClothesID      utils.UUID
	ClothesDir     string
	Classification ClothesClassificationModelResponse
}

//easyjson:json
type ClothesClassificationRequest struct { // Request to ML-server
	Tags          []string `json:"tags,!omitempty"` //lint:ignore SA5008 easyjson custom tags
	Styles        []string
	Categories    []string
	Subcategories []string
	Seasons       []string
}

//easyjson:json
type ClothesClassificationResponse struct { // End-user response
	Type    utils.UUID
	Subtype utils.UUID
	Style   utils.UUID
	Seasons []Season `json:"seasons,!omitempty"` //lint:ignore SA5008 easyjson custom tags
	Tags    []string `json:"tags,!omitempty"`    //lint:ignore SA5008 easyjson custom tags
}

//easyjson:json
type ClothesClassificationModelResponse struct {
	Tags          map[string]float32
	Categories    map[string]float32
	Subcategories map[string]float32
	Seasons       map[string]float32
	Styles        map[string]float32
}
