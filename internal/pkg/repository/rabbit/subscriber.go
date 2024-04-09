package rabbit

import (
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type Subscriber[T any, PT interface {
	*T
	easyjson.Unmarshaler
}] struct {
	rabbit *rabbitmq.Conn
	queue  string
}

func NewSubscriber[T any, PT interface {
	*T
	easyjson.Unmarshaler
}](rabbit *rabbitmq.Conn, queue string) domain.Subscriber[T] {
	return &Subscriber[T, PT]{
		rabbit: rabbit,
		queue:  queue,
	}
}

func (s *Subscriber[T, PT]) Listen(logger *zap.SugaredLogger, handler func(*T) domain.Result) error {
	consumer, err := rabbitmq.NewConsumer(
		s.rabbit,
		s.queue,
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

		logger.Infow("rabbit", "queue", s.queue, "got", string(delivery.Body))

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
