package ml

import (
	"context"
	"time"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClothesProcessor struct {
	ch    *amqp.Channel
	queue amqp.Queue
}

func New(queueName string, ch *amqp.Channel) (domain.ClothesProcessingModel, error) {
	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &ClothesProcessor{
		ch:    ch,
		queue: queue,
	}, nil
}

func (p *ClothesProcessor) Process(ctx context.Context, opts domain.ClothesProcessingOpts) error {
	bytes, err := easyjson.Marshal(opts)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"",
		p.queue.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  common.ContentTypeJSON,
			Body:         bytes,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func (p *ClothesProcessor) GetTryOnResults() (chan interface{}, error) {
	// p.ch.Consume(p.queue.Name, )
	return nil, nil
}

func (p *ClothesProcessor) TryOn(ctx context.Context, opts domain.TryOnOpts) error {
	bytes, err := easyjson.Marshal(opts)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"",
		p.queue.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  common.ContentTypeJSON,
			Body:         bytes,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}
