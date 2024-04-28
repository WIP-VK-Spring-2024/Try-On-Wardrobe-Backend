package outfitgen

import (
	"context"
	"fmt"
	"strings"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/clothes"
	"try-on/internal/pkg/repository/sqlc/outfits"
	"try-on/internal/pkg/utils"

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
	translator domain.Translator,
) domain.OutfitGenerator {
	return &OutfitGenerator{
		publisher:  publisher,
		subscriber: subscriber,
		clothes:    clothes.New(db),
		outfits:    outfits.New(db),
		weather:    weather,
		translator: translator,
	}
}

func (gen *OutfitGenerator) Close() {
	gen.publisher.Close()
}

func (gen *OutfitGenerator) Generate(ctx context.Context, request domain.OutfitGenerationRequest) error {
	var tempCelcius *int = nil

	if request.UseWeather {
		weather, err := gen.weather.CurrentWeather(request.Pos)
		if err != nil {
			return err
		}
		fmt.Println("Got weather: ", weather.Temp)

		tmp := int(weather.Temp)
		tempCelcius = &tmp
	}

	clothes, err := gen.clothes.GetByWeather(request.UserID, tempCelcius)
	if err != nil {
		return err
	}

	maxAmount := maxOutfitNum(clothes)
	if maxAmount == 0 {
		return app_errors.ErrNotEnoughClothes
	}

	if request.Amount > maxAmount {
		request.Amount = maxAmount
	}

	purposes, err := gen.outfits.GetPurposeEngNames(request.Purposes)
	if err != nil {
		return err
	}

	translatedPrompt := ""

	if request.Prompt != "" {
		translatedPrompt, err = gen.translator.Translate(request.Prompt, domain.LanguageRU, domain.LanguageEN)
		if err != nil {
			return err
		}
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
	ctx := middleware.WithLogger(context.Background(), logger)

	return gen.subscriber.Listen(ctx, handler)
}

func maxOutfitNum(clothes []domain.GenClothesInfo) int {
	upperNum := utils.Count(clothes, func(elem domain.GenClothesInfo) bool {
		return elem.Category == domain.GenCategoryUpper
	})

	lowerNum := utils.Count(clothes, func(elem domain.GenClothesInfo) bool {
		return elem.Category == domain.GenCategoryLower
	})

	outerNum := utils.Count(clothes, func(elem domain.GenClothesInfo) bool {
		return elem.Category == domain.GenCategoryOuter
	})

	dressNum := utils.Count(clothes, func(elem domain.GenClothesInfo) bool {
		return elem.Category == domain.GenCategoryDress
	})

	return upperNum*lowerNum*(outerNum+1) + dressNum
}
