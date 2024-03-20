package domain

import (
	"context"
	"io"
)

type FileManager interface {
	Save(ctx context.Context, dir, name string, data io.Reader) error
	Delete(ctx context.Context, dir, name string) error
}
