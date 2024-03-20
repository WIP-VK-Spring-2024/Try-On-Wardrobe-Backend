package ml

import (
	"context"
	"time"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type ClothesProcessor struct {
	publisher     *rabbitmq.Publisher
	consumer      *rabbitmq.Consumer
	requestQueue  string
	responseQueue string
}

func (p *ClothesProcessor) Close() {
	p.publisher.Close()
}

func New(
	requestQueue string,
	responseQueue string,
	rabbit *rabbitmq.Conn,
) (domain.ClothesProcessingModel, error) {
	publisher, err := rabbitmq.NewPublisher(
		rabbit,
		rabbitmq.WithPublisherOptionsExchangeName(requestQueue),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(
		rabbit,
		responseQueue,
		rabbitmq.WithConsumerOptionsExchangeName(responseQueue),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		publisher.Close()
		return nil, err
	}

	return &ClothesProcessor{
		publisher:     publisher,
		consumer:      consumer,
		requestQueue:  requestQueue,
		responseQueue: responseQueue,
	}, nil
}

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingOpts) error {
	return p.publish(ctx, opts, p.requestQueue)
}

func (p *ClothesProcessor) GetTryOnResults(logger *zap.SugaredLogger, handler func(*domain.TryOnResponse) domain.Result) error {
	defer p.consumer.Close()

	return p.consumer.Run(func(delivery rabbitmq.Delivery) rabbitmq.Action {
		logger.Infow("rabbit", "got", string(delivery.Body))

		var resp domain.TryOnResponse
		err := easyjson.Unmarshal(delivery.Body, &resp)
		if err != nil {
			logger.Infow("rabbit", "error", err)
			return rabbitmq.NackDiscard
		}

		res := handler(&resp)
		switch res {
		case domain.ResultOk:
			return rabbitmq.Ack

		case domain.ResultRetry:
			return rabbitmq.NackRequeue

		case domain.ResultDiscard:
			fallthrough

		default:
			return rabbitmq.NackDiscard
		}
	})
}

func (p *ClothesProcessor) TryOn(ctx context.Context, opts domain.TryOnOpts) error {
	return p.publish(ctx, opts, p.requestQueue)
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
