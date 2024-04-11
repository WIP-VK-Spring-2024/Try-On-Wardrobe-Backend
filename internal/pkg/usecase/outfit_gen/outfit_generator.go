package outfitgen

import (
	"context"

	"try-on/internal/pkg/domain"

	"go.uber.org/zap"
)

type OutfitGenerator struct {
	publisher  domain.Publisher[domain.OutfitGenerationModelRequest]
	subscriber domain.Subscriber[domain.OutfitGenerationResponse]
	clothes    domain.ClothesRepository
	outfits    domain.OutfitRepository
}

func New(
	publisher domain.Publisher[domain.OutfitGenerationModelRequest],
	subscriber domain.Subscriber[domain.OutfitGenerationResponse],
	clothes domain.ClothesRepository,
	outfits domain.OutfitRepository,
) domain.OutfitGenerator {
	return &OutfitGenerator{
		publisher:  publisher,
		subscriber: subscriber,
		clothes:    clothes,
		outfits:    outfits,
	}
}

func (gen *OutfitGenerator) Generate(ctx context.Context, request domain.OutfitGenerationRequest) error {
	// Get current weather
	// Get clothes matching the weather
	// Translate tags and prompt
	// Concatenate translated tags with prompt
	// return gen.publisher.Publish(context.Background(), )
	return nil
}

func (gen *OutfitGenerator) ListenGenerationResults(logger *zap.SugaredLogger, handler func(*domain.OutfitGenerationResponse) domain.Result) error {
	// handler sends result to centrifugo
	return gen.subscriber.Listen(logger, handler)
}
