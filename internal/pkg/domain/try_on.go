package domain

import (
	"context"

	"try-on/internal/pkg/utils"

	"go.uber.org/zap"
)

//easyjson:json
type TryOnResult struct {
	Model

	Image       string
	Rating      int `json:"rating,!omitempty"` //lint:ignore SA5008 easyjson custom tags
	UserImageID utils.UUID
	ClothesID   []utils.UUID
	OutfitID    utils.UUID
}

type TryOnUsecase interface {
	Closer
	TryOn(ctx context.Context, clothes []utils.UUID, opts TryOnOpts) error
	TryOnOutfit(ctx context.Context, outfit utils.UUID, opts TryOnOpts) error
	GetTryOnResults(logger *zap.SugaredLogger, handler func(*TryOnResponse) Result) error
}

//easyjson:json
type TryOnOpts struct {
	UserID       utils.UUID
	UserImageID  utils.UUID
	UserImageDir string
	ClothesDir   string
}

//easyjson:json
type TryOnRequest struct {
	TryOnOpts
	OutfitID utils.UUID
	Clothes  []TryOnClothesInfo
}

//easyjson:json
type TryOnClothesInfo struct {
	ClothesID utils.UUID
	Category  string
	Layer     int `json:"-"`
}

//easyjson:json
type TryOnResponse struct {
	QueueResponse

	UserID      utils.UUID
	OutfitID    utils.UUID
	Clothes     []TryOnClothesInfo
	UserImageID utils.UUID
	TryOnID     string
	TryOnDir    string
}

type TryOnResultRepository interface {
	Create(res *TryOnResult) error
	Delete(id utils.UUID) error
	SetTryOnResultID(outfitId, id utils.UUID) error
	GetByUser(userID utils.UUID) ([]TryOnResult, error)
	Get(id utils.UUID) (*TryOnResult, error)
	GetByClothes(clothesID utils.UUID) ([]TryOnResult, error)
	Rate(id utils.UUID, rating int) error
}
