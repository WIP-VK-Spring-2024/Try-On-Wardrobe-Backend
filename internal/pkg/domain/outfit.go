package domain

import (
	"context"
	"time"

	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"

	"go.uber.org/zap"
)

type Privacy string

const (
	PrivacyPublic  Privacy = "public"
	PrivacyPrivate Privacy = "private"
	PrivacyFriends Privacy = "public"
)

//easyjson:json
type Outfit struct {
	Model

	UserID    utils.UUID
	StyleID   utils.UUID
	PurposeID utils.UUID

	Privacy Privacy

	Name       string
	Note       optional.String
	Image      string
	Transforms TransformMap
	Seasons    []Season
	Tags       []string
}

//easyjson:json
type OutfitPurpose struct {
	Model

	Name    string
	EngName string
}

//easyjson:json
type TransformMap map[utils.UUID]Transform

//easyjson:json
type Transform struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
	Angle  float32
	Scale  float32
}

//easyjson:json
type GeoPosition struct {
	Lat float32
	Lon float32
}

//easyjson:json
type OutfitGenerationRequest struct {
	UserID   utils.UUID
	IP       string
	Pos      WeatherRequest
	Purposes []string
	Prompt   string
}

//easyjson:json
type OutfitGenerationModelRequest struct {
	UserID       utils.UUID
	Clothes      []GenClothesInfo
	Prompt       string
	SampleAmount int
}

//easyjson:json
type GenClothesInfo struct {
	ClothesID utils.UUID
	Category  string
}

//easyjson:json
type OutfitGenerationResponse struct {
	UserID  utils.UUID
	Outfits [][]utils.UUID
}

type OutfitGenerator interface {
	Generate(ctx context.Context, request OutfitGenerationRequest) error
	ListenGenerationResults(logger *zap.SugaredLogger, handler func(*OutfitGenerationResponse) Result) error
}

type OutfitUsecase interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(userId, outfitID utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) ([]TryOnClothesInfo, error)
	GetByUser(utils.UUID) ([]Outfit, error)
	GetOutfitPurposes() ([]OutfitPurpose, error)
}

type OutfitRepository interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) ([]TryOnClothesInfo, error)
	GetByUser(utils.UUID) ([]Outfit, error)
	GetOutfitPurposesByEngName(engNames []string) ([]OutfitPurpose, error)
	GetOutfitPurposes() ([]OutfitPurpose, error)
	GetPurposeEngNames(names []string) ([]string, error)
}
