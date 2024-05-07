package domain

import (
	"context"
)

type Closer interface {
	Close()
}

type Publisher[T any] interface {
	Closer
	Publish(ctx context.Context, request T) error
}

type ChannelPublisher[T any] interface {
	Publish(ctx context.Context, channel string, message T) error
}

type Subscriber[T any] interface {
	Listen(ctx context.Context, handler func(response *T) Result) error
}

//easyjson:json
type QueueResponse struct {
	StatusCode int
	Message    string
}
