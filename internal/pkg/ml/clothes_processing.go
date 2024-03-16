package ml

import (
	"context"
	"log"
	"time"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClothesProcessor struct {
	ch           *amqp.Channel
	queue        amqp.Queue
	reponseQueue amqp.Queue
}

func New(queueName string, reponseQueueName string, ch *amqp.Channel) (domain.ClothesProcessingModel, error) {
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

	reponseQueue, err := ch.QueueDeclare(
		reponseQueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

	return &ClothesProcessor{
		ch:           ch,
		queue:        queue,
		reponseQueue: reponseQueue,
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

func (p *ClothesProcessor) GetTryOnResults() (chan domain.TryOnResponse, error) {
	ch, err := p.ch.Consume(
		p.reponseQueue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	respChan := make(chan domain.TryOnResponse, 3)

	var resp domain.TryOnResponse
	go func() {
		for delivery := range ch {
			err := easyjson.Unmarshal(delivery.Body, &resp)
			if err != nil {
				log.Println("ERROR:", err)
			}
			respChan <- resp
		}
	}()

	return respChan, nil
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
