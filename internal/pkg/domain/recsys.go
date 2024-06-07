package domain

import (
	"context"

	"try-on/internal/pkg/utils"

	"go.uber.org/zap"
)

//easyjson:json
type RecsysRequest struct {
	UserID        utils.UUID
	SamplesAmount int
}

//easyjson:json
type RecsysResponse struct {
	QueueResponse

	UserID    utils.UUID
	OutfitIds []utils.UUID
}

type Recsys interface {
	Closer
	GetRecommendations(ctx context.Context, limit int, request RecsysRequest) ([]Post, error)
	ListenResults(logger *zap.SugaredLogger) error
}
