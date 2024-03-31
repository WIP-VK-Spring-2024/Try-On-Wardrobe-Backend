package ml

import (
	"context"
	"time"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/subtypes"
	"try-on/internal/pkg/repository/sqlc/types"
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

	types    domain.TypeRepository
	subtypes domain.SubtypeRepository
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
		publisher: publisher,
		rabbit:    rabbit,
		tryOn:     tryOn,
		process:   process,
		cfg:       cfg,
		types:     types.New(db),
		subtypes:  subtypes.New(db),
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
	Classification classification
}

//easyjson:json
type classification struct {
	Tags          map[string]float32
	Categories    map[string]float32
	Subcategories map[string]float32
	Seasons       map[string]float32
}

func (p *ClothesProcessor) notPassesThreshold(_ string, value float32) bool {
	return value < p.cfg.Threshold
}

func (p *ClothesProcessor) GetProcessingResults(logger *zap.SugaredLogger, handler func(*domain.ClothesProcessingResponse) domain.Result) error {
	return getResults(p, p.process, logger, func(result *processingResult) domain.Result {
		maps.DeleteFunc(result.Classification.Tags, p.notPassesThreshold)

		maps.DeleteFunc(result.Classification.Seasons, p.notPassesThreshold)

		maps.DeleteFunc(result.Classification.Subcategories, p.notPassesThreshold)

		// TODO: Get type/subtype ids based on ML classifiers from repos

		return handler(&domain.ClothesProcessingResponse{
			UserID:    result.UserID,
			ClothesID: result.ClothesID,
			ResultDir: result.ResultDir,
			Classification: domain.ClothesClassificationResponse{
				Tags:    utils.SortedKeysByValue(result.Classification.Tags),
				Seasons: maps.Keys(result.Classification.Seasons),
			},
		})
	})
}

func (p *ClothesProcessor) TryOn(ctx context.Context, opts domain.TryOnRequest) error {
	return p.publish(ctx, opts, p.tryOn.Request)
}

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingRequest) error {
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
