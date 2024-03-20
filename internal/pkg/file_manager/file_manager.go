package file_manager

import (
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/file_manager/filesystem"
	httpuploader "try-on/internal/pkg/file_manager/http_uploader"
	"try-on/internal/pkg/file_manager/s3"
)

const (
	Filesystem = "fs"
	S3         = "s3"
	Http       = "http"
)

func New(cfg *config.Static) (domain.FileManager, error) {
	switch cfg.Type {
	case Http:
		return httpuploader.New(&cfg.HttpApi), nil

	case S3:
		return s3.New(&cfg.S3)

	case Filesystem:
		fallthrough

	default:
		return filesystem.New(cfg.Dir), nil
	}
}
