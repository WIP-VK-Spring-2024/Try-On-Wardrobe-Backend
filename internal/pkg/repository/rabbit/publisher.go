package rabbit

import (
	"context"
	"log"
	"time"

	"try-on/internal/pkg/common"

	"github.com/mailru/easyjson"
	"github.com/wagslane/go-rabbitmq"
)

type Publisher[T easyjson.Marshaler] struct {
	publisher *rabbitmq.Publisher
	queue     string
}

func NewPublisher[T easyjson.Marshaler](conn *rabbitmq.Conn, queue string) *Publisher[T] {
	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		panic(err)
	}

	return &Publisher[T]{
		publisher: publisher,
		queue:     queue,
	}
}

func (p Publisher[_]) Close() {
	p.publisher.Close()
}

func (p Publisher[T]) Publish(ctx context.Context, req T) error {
	bytes, err := easyjson.Marshal(req)
	if err != nil {
		return err
	}

	log.Printf("Publishing to rabbit: queue is '%s', payload is '%s'\n", p.queue, string(bytes))

	return p.publisher.PublishWithContext(
		ctx,
		bytes,
		[]string{p.queue},
		rabbitmq.WithPublishOptionsContentType(common.ContentTypeJSON),
		rabbitmq.WithPublishOptionsTimestamp(time.Now()),
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}
