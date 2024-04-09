package ml

import (
	"cmp"
	"context"
	"strings"

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

func (p *ClothesProcessor) GetProcessingResults(logger *zap.SugaredLogger, handler func(*domain.ClothesProcessingResponse) domain.Result) error {
	return p.subscriber.Listen(logger, func(result *domain.ClothesProcessingModelResponse) domain.Result {
		maps.DeleteFunc(result.Classification.Tags, notPassesThreshold[string](p.cfg.Threshold))

		maps.DeleteFunc(result.Classification.Seasons, notPassesThreshold[string](p.cfg.Threshold))

		styleId, err := p.classificationRepo.GetStyleId(maxKey(result.Classification.Styles))
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		typeId, err := p.classificationRepo.GetTypeId(maxKey(result.Classification.Categories))
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		subcategories := filterSubcategories(result.Classification.Subcategories, p.cfg.Threshold)
		subtypeIds, err := p.classificationRepo.GetSubtypeIds(subcategories)
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		tags, err := p.classificationRepo.GetTags(utils.SortedKeysByValue(result.Classification.Tags))
		if err != nil {
			logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		return handler(&domain.ClothesProcessingResponse{
			UserID:     result.UserID,
			ClothesID:  result.ClothesID,
			ClothesDir: result.ClothesDir,
			Classification: domain.ClothesClassificationResponse{
				Tags:     tags,
				Seasons:  removeClothesSuffix(maps.Keys(result.Classification.Seasons)),
				Style:    styleId,
				Type:     typeId,
				Subtypes: subtypeIds,
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

func filterSubcategories(subcategories map[string]float32, threshold float32) []string {
	tmp := maps.Clone(subcategories)
	maps.DeleteFunc(tmp, notPassesThreshold[string](threshold))

	if len(tmp) != 0 {
		return maps.Keys(tmp)
	}

	sorted := utils.SortedKeysByValue(subcategories)
	return sorted[:len(sorted)/2]
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
