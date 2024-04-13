package outfitgen

import (
	"context"
	"fmt"
	"strings"

	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/clothes"
	"try-on/internal/pkg/repository/sqlc/outfits"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
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
	db *pgxpool.Pool,
	weather domain.WeatherService,
) domain.OutfitGenerator {
	return &OutfitGenerator{
		publisher:  publisher,
		subscriber: subscriber,
		clothes:    clothes.New(db),
		outfits:    outfits.New(db),
		weather:    weather,
	}
}

func (gen *OutfitGenerator) Close() {
	gen.publisher.Close()
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

	purposes, err := gen.outfits.GetPurposeEngNames(request.Purposes)
	if err != nil {
		return err
	}

	translatedPrompt, err := gen.translator.Translate(request.Prompt, domain.LanguageRU, domain.LanguageEN)
	if err != nil {
		return err
	}

	purposes = append([]string{translatedPrompt}, purposes...)

	modelRequest := domain.OutfitGenerationModelRequest{
		UserID:  request.UserID,
		Clothes: clothes,
		Amount:  request.Amount,
		Prompt:  strings.Join(purposes, ". "),
	}

	bytes, _ := easyjson.Marshal(modelRequest)
	fmt.Println(string(bytes))

	return gen.publisher.Publish(ctx, modelRequest)
}

func (gen *OutfitGenerator) ListenGenerationResults(logger *zap.SugaredLogger, handler func(*domain.OutfitGenerationResponse) domain.Result) error {
	// TODO:Post-processing?
	return gen.subscriber.Listen(logger, handler)
}