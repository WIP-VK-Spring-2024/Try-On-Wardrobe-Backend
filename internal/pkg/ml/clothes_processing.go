package ml

import (
	"context"
	"time"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type ClothesProcessor struct {
	publisher *rabbitmq.Publisher
	rabbit    *rabbitmq.Conn
	tryOn     config.RabbitQueue
	process   config.RabbitQueue
}

func (p *ClothesProcessor) Close() {
	p.publisher.Close()
}

func New(
	tryOn config.RabbitQueue,
	process config.RabbitQueue,
	rabbit *rabbitmq.Conn,
) (domain.ClothesProcessingModel, error) {
	publisher, err := rabbitmq.NewPublisher(
		rabbit,
		rabbitmq.WithPublisherOptionsExchangeName("default_exchange"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &ClothesProcessor{
		publisher: publisher,
		rabbit:    rabbit,
		tryOn:     tryOn,
		process:   process,
	}, nil
}

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingRequest) error {
	return p.publish(ctx, opts, p.process.Request)
}

func (p *ClothesProcessor) GetTryOnResults(logger *zap.SugaredLogger, handler func(*domain.TryOnResponse) domain.Result) error {
	consumer, err := rabbitmq.NewConsumer(
		p.rabbit,
		p.tryOn.Response,
		rabbitmq.WithConsumerOptionsExchangeName(p.tryOn.Response),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}
	defer consumer.Close()

	return consumer.Run(func(delivery rabbitmq.Delivery) rabbitmq.Action {
		logger.Infow("rabbit", "got", string(delivery.Body))

		var resp domain.TryOnResponse
		err := easyjson.Unmarshal(delivery.Body, &resp)
		if err != nil {
			logger.Infow("rabbit", "error", err)
			return rabbitmq.NackDiscard
		}

		return toRabbitAction(handler(&resp))
	})
}

func (p *ClothesProcessor) GetProcessingResults(logger *zap.SugaredLogger, handler func(*domain.ClothesProcessingResponse) domain.Result) error {
	consumer, err := rabbitmq.NewConsumer(
		p.rabbit,
		p.process.Response,
		rabbitmq.WithConsumerOptionsExchangeName(p.process.Response),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}
	defer consumer.Close()

	return consumer.Run(func(delivery rabbitmq.Delivery) rabbitmq.Action {
		logger.Infow("rabbit", "got", string(delivery.Body))

		var resp domain.ClothesProcessingResponse
		err := easyjson.Unmarshal(delivery.Body, &resp)
		if err != nil {
			logger.Infow("rabbit", "error", err)
			return rabbitmq.NackDiscard
		}

		return toRabbitAction(handler(&resp))
	})
}

func (p *ClothesProcessor) TryOn(ctx context.Context, opts domain.TryOnRequest) error {
	return p.publish(ctx, opts, p.tryOn.Request)
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
