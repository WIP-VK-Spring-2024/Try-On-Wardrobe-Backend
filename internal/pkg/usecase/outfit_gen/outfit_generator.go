package outfitgen

import (
	"context"

	"try-on/internal/pkg/domain"

	"go.uber.org/zap"
)

type OutfitGenerator struct {
	publisher  domain.Publisher[domain.OutfitGenerationModelRequest]
	subscriber domain.Subscriber[domain.OutfitGenerationResponse]

	clothes domain.ClothesRepository
	outfits domain.OutfitRepository

	weather    domain.WeatherService
	translator domain.Translator
}

func New(
	publisher domain.Publisher[domain.OutfitGenerationModelRequest],
	subscriber domain.Subscriber[domain.OutfitGenerationResponse],
	clothes domain.ClothesRepository,
	outfits domain.OutfitRepository,
	weather domain.WeatherService,
) domain.OutfitGenerator {
	return &OutfitGenerator{
		publisher:  publisher,
		subscriber: subscriber,
		clothes:    clothes,
		outfits:    outfits,
		weather:    weather,
	}
}

func (gen *OutfitGenerator) Generate(ctx context.Context, request domain.OutfitGenerationRequest) error {
	weather, err := gen.weather.CurrentWeather(request.Pos)
	if err != nil {
		return err
	}

	clothes, err := gen.clothes.GetByWeather(request.UserID, int(weather.Temp))
	if err != nil {
		return err
	}

	// TODO: Get tags

	translatedPrompt, err := gen.translator.Translate(request.Prompt, domain.LanguageRU, domain.LanguageEN)
	if err != nil {
		return err
	}

	// Translate tags and prompt

	return gen.publisher.Publish(context.Background(), domain.OutfitGenerationModelRequest{
		UserID:  request.UserID,
		Clothes: clothes,
		Prompt:  translatedPrompt, /* + genTags */
	})
}

func (gen *OutfitGenerator) ListenGenerationResults(logger *zap.SugaredLogger, handler func(*domain.OutfitGenerationResponse) domain.Result) error {
	// handler sends result to centrifugo
	return gen.subscriber.Listen(logger, handler)
}
