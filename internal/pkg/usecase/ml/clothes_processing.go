package ml

import (
	"cmp"
	"context"
	"slices"
	"time"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/classification"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type ClothesProcessor struct {
	publisher *rabbitmq.Publisher
	rabbit    *rabbitmq.Conn

	tryOn   config.RabbitQueue
	process config.RabbitQueue

	cfg *config.Classification

	classificationRepo domain.ClothesClassificationRepository
}

func (p *ClothesProcessor) Close() {
	p.publisher.Close()
}

func New(
	tryOn config.RabbitQueue,
	process config.RabbitQueue,
	rabbit *rabbitmq.Conn,
	cfg *config.Classification,
	db *pgxpool.Pool,
) (domain.ClothesProcessingModel, error) {
	publisher, err := rabbitmq.NewPublisher(
		rabbit,
	)
	if err != nil {
		return nil, err
	}

	return &ClothesProcessor{
		publisher:          publisher,
		rabbit:             rabbit,
		tryOn:              tryOn,
		process:            process,
		cfg:                cfg,
		classificationRepo: classification.New(db),
	}, nil
}

type handlerFunc[T easyjson.Unmarshaler] func(T) domain.Result

func (p *ClothesProcessor) GetTryOnResults(logger *zap.SugaredLogger, handler func(*domain.TryOnResponse) domain.Result) error {
	return getResults(p, p.tryOn, logger, handler)
}

//easyjson:json
type processingResult struct {
	UserID         utils.UUID
	ClothesID      utils.UUID
	ResultDir      string
	Classification classificationModelResponse
}

//easyjson:json
type classificationModelResponse struct {
	Tags          map[string]float32
	Categories    map[string]float32
	Subcategories map[string]float32
	Seasons       map[domain.Season]float32
	Styles        map[string]float32
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

func (p *ClothesProcessor) GetProcessingResults(logger *zap.SugaredLogger, handler func(*domain.ClothesProcessingResponse) domain.Result) error {
	return getResults(p, p.process, logger, func(result *processingResult) domain.Result {
		maps.DeleteFunc(result.Classification.Tags, notPassesThreshold[string](p.cfg.Threshold))

		maps.DeleteFunc(result.Classification.Seasons, notPassesThreshold[domain.Season](p.cfg.Threshold))

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
			ClothesDir: result.ResultDir,
			Classification: domain.ClothesClassificationResponse{
				Tags:     tags,
				Seasons:  maps.Keys(result.Classification.Seasons),
				Style:    styleId,
				Type:     typeId,
				Subtypes: subtypeIds,
			},
		})
	})
}

func (p *ClothesProcessor) TryOn(ctx context.Context, opts domain.TryOnRequest) error {
	if err := validateTryOnCategories(opts.Clothes); err != nil {
		return err
	}
	return p.publish(ctx, opts, p.tryOn.Request)
}

func validateTryOnCategories(clothes []domain.TryOnClothesInfo) error {
	if len(clothes) < 1 || len(clothes) > 2 {
		return app_errors.ErrTryOnInvalidClothesNum
	}

	dressIdx := slices.IndexFunc(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryDress
	})
	if dressIdx != -1 && len(clothes) != 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	upperCount := utils.Count(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryUpper
	})
	if upperCount > 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	lowerCount := utils.Count(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryLower
	})
	if lowerCount > 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	return nil
}

const tagLimit = 10

const clothesSuffix = " clothes"

func addClothesSuffix(seasons []string) []string {
	result := make([]string, 0, len(seasons))

	for _, season := range seasons {
		result = append(result, season+clothesSuffix)
	}
	return result
}

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingRequest) error {
	classificationRequest, err := p.classificationRepo.GetClassifications(opts.UserID, tagLimit)
	if err != nil {
		return err
	}

	classificationRequest.Seasons = addClothesSuffix(classificationRequest.Seasons)

	opts.Classification = *classificationRequest
	return p.publish(ctx, opts, p.process.Request)
}

func (p *ClothesProcessor) publish(ctx context.Context, payload easyjson.Marshaler, routingKeys ...string) error {
	bytes, err := easyjson.Marshal(payload)
	if err != nil {
		return err
	}

	return p.publisher.PublishWithContext(
		ctx,
		bytes,
		routingKeys,
		rabbitmq.WithPublishOptionsContentType(common.ContentTypeJSON),
		rabbitmq.WithPublishOptionsTimestamp(time.Now()),
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}

func getResults[T any, PT interface {
	*T
	UnmarshalEasyJSON(w *jlexer.Lexer)
}](p *ClothesProcessor, queue config.RabbitQueue, logger *zap.SugaredLogger, handler handlerFunc[PT]) error {
	consumer, err := rabbitmq.NewConsumer(
		p.rabbit,
		queue.Response,
	)
	if err != nil {
		return err
	}
	defer consumer.Close()

	return consumer.Run(func(delivery rabbitmq.Delivery) rabbitmq.Action {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorln(err)
			}
		}()

		logger.Infow("rabbit", "got", string(delivery.Body))

		resp := PT(new(T))

		err := easyjson.Unmarshal(delivery.Body, resp)
		if err != nil {
			logger.Infow("rabbit", "error", err)
			return rabbitmq.NackDiscard
		}

		return toRabbitAction(handler(resp))
	})
}

func toRabbitAction(result domain.Result) rabbitmq.Action {
	switch result {
	case domain.ResultOk:
		return rabbitmq.Ack

	case domain.ResultRetry:
		return rabbitmq.NackRequeue

	case domain.ResultDiscard:
		fallthrough

	default:
		return rabbitmq.NackDiscard
	}
}
