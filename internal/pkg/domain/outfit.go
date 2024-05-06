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
	PrivacyFriends Privacy = "friends"

	GenCategoryUpper = "upper garment"
	GenCategoryLower = "lower garment"
	GenCategoryOuter = "outerwear"
	GenCategoryDress = "dress"
)

var Privacies = map[Privacy]struct{}{
	PrivacyPublic: {}, PrivacyPrivate: {}, PrivacyFriends: {},
}

//easyjson:json
type Outfit struct {
	Model

	UserID        utils.UUID
	StyleID       utils.UUID
	PurposeID     utils.UUID
	TryOnResultID utils.UUID

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
	ZIndex int `json:"z_index,!omitempty"` //lint:ignore SA5008 easyjson custom tags
}

//easyjson:json
type GeoPosition struct {
	Lat float32 `query:"lat"`
	Lon float32 `query:"lon"`
}

//easyjson:json
type OutfitGenerationRequest struct {
	UserID     utils.UUID
	Amount     int            `query:"amount"`
	UseWeather bool           `query:"use_weather"`
	Pos        WeatherRequest `query:"pos"`
	Purposes   []string       `query:"purposes"`
	Prompt     string         `query:"prompt"`
}

//easyjson:json
type OutfitGenerationModelRequest struct {
	UserID  utils.UUID
	Clothes []GenClothesInfo
	Prompt  string
	Amount  int
}

//easyjson:json
type GenClothesInfo struct {
	ClothesID utils.UUID
	Category  string
}

//easyjson:json
type OutfitGenClothes struct {
	ClothesID utils.UUID
}

//easyjson:json
type OutfitGenOutfit struct {
	Clothes []OutfitGenClothes
}

//easyjson:json
type OutfitGenerationResponse struct {
	QueueResponse

	UserID  utils.UUID
	Outfits []OutfitGenOutfit
}

type OutfitGenerator interface {
	Closer
	Generate(ctx context.Context, request OutfitGenerationRequest) error
	ListenGenerationResults(logger *zap.SugaredLogger, handler func(*OutfitGenerationResponse) Result) error
}

type OutfitUsecase interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(userId, outfitId utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) ([]TryOnClothesInfo, error)
	GetByUser(userId utils.UUID, publicOnly bool) ([]Outfit, error)
	GetOutfitPurposes() ([]OutfitPurpose, error)
}

type OutfitRepository interface {
	Create(*Outfit) error
	Update(*Outfit) error
	Delete(utils.UUID) error
	GetById(utils.UUID) (*Outfit, error)
	Get(since time.Time, limit int) ([]Outfit, error)
	GetClothesInfo(utils.UUID) ([]TryOnClothesInfo, error)
	GetByUser(userId utils.UUID, publicOnly bool) ([]Outfit, error)
	GetOutfitPurposesByEngName(engNames []string) ([]OutfitPurpose, error)
	GetOutfitPurposes() ([]OutfitPurpose, error)
	GetPurposeEngNames(names []string) ([]string, error)
}
