package ml

import (
	"cmp"
	"context"
	"log"
	"strings"

	"try-on/internal/middleware"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/classification"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type ClothesProcessor struct {
	publisher  domain.Publisher[domain.ClothesProcessingRequest]
	subscriber domain.Subscriber[domain.ClothesProcessingModelResponse]

	cfg                *config.Classification
	classificationRepo domain.ClothesClassificationRepository
}

func (p *ClothesProcessor) Close() {
	p.publisher.Close()
}

func New(
	pub domain.Publisher[domain.ClothesProcessingRequest],
	sub domain.Subscriber[domain.ClothesProcessingModelResponse],
	cfg *config.Classification,
	db *pgxpool.Pool,
) domain.ClothesProcessingModel {
	return &ClothesProcessor{
		cfg:                cfg,
		classificationRepo: classification.New(db),
		publisher:          pub,
		subscriber:         sub,
	}
}

const tagLimit = 10

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingRequest) error {
	classificationRequest, err := p.classificationRepo.GetClassifications(opts.UserID, tagLimit)
	if err != nil {
		return err
	}

	classificationRequest.Seasons = addClothesSuffix(classificationRequest.Seasons)

	opts.Classification = *classificationRequest
	return p.publisher.Publish(ctx, opts)
}

func maxKeys(input map[string]float32, threshold float32) []string {
	maps.DeleteFunc(input, notPassesThreshold[string](threshold))
	return maps.Keys(input)
}

func (p *ClothesProcessor) GetProcessingResults(logger *zap.SugaredLogger, handler func(*domain.ClothesProcessingResponse) domain.Result) error {
	ctx := middleware.WithLogger(context.Background(), logger)

	return p.subscriber.Listen(ctx, func(result *domain.ClothesProcessingModelResponse) domain.Result {
		maps.DeleteFunc(result.Classification.Tags, notPassesThreshold[string](p.cfg.Threshold))

		log.Println("Filtered tags:", maxKeys(result.Classification.Tags, p.cfg.Threshold))
		log.Println("Filtered seasons:", maxKeys(result.Classification.Seasons, p.cfg.Threshold))

		maxStyle := maxKey(result.Classification.Styles)
		styleId, err := p.classificationRepo.GetStyleId(maxStyle)
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		subcategory := maxKey(result.Classification.Subcategories)
		log.Println("Filtered subcategories:", subcategory)
		subtypeId, err := p.classificationRepo.GetSubtypeIds(subcategory)
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		tags, err := p.classificationRepo.GetTags(utils.SortedKeysByValue(result.Classification.Tags))
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		typeId, tryonable, err := p.classificationRepo.GetTypeBySubtype(subtypeId)
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		return handler(&domain.ClothesProcessingResponse{
			UserID:     result.UserID,
			ClothesID:  result.ClothesID,
			ClothesDir: result.ClothesDir,
			Tryonable:  tryonable,
			Classification: domain.ClothesClassificationResponse{
				Tags:    tags,
				Seasons: removeClothesSuffix(maxKeys(result.Classification.Seasons, p.cfg.Threshold)),
				Style:   styleId,
				Type:    typeId,
				Subtype: subtypeId,
			},
		})
	})
}

func notPassesThreshold[T ~string](threshold float32) func(_ T, value float32) bool {
	return func(_ T, value float32) bool {
		return value < threshold
	}
}

func maxKey[M ~map[K]V, K comparable, V cmp.Ordered](input M) K {
	var curMax V
	var result K

	for key, value := range input {
		if value > curMax {
			curMax = value
			result = key
		}
	}

	return result
}

const clothesSuffix = " clothes"

func addClothesSuffix(seasons []string) []string {
	result := make([]string, 0, len(seasons))

	for _, season := range seasons {
		result = append(result, season+clothesSuffix)
	}
	return result
}

func removeClothesSuffix(seasons []string) []domain.Season {
	result := make([]domain.Season, 0, len(seasons))

	for _, season := range seasons {
		cut, _ := strings.CutSuffix(season, clothesSuffix)
		result = append(result, domain.Season(cut))
	}
	return result
}
