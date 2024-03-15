package domain

import (
	"context"
	"io"
)

type FileManager interface {
	Save(ctx context.Context, dir, name string, data io.Reader) error
	Get(ctx context.Context, dir, name string) (io.ReadCloser, error)
	Delete(ctx context.Context, dir, name string) error
}
