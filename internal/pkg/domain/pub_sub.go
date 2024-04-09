package domain

import (
	"context"

	"go.uber.org/zap"
)

type Closer interface {
	Close()
}

type Publisher[T any] interface {
	Closer
	Publish(ctx context.Context, request T) error
}

type Subscriber[T any] interface {
	Listen(logger *zap.SugaredLogger, handler func(*T) Result) error
}
