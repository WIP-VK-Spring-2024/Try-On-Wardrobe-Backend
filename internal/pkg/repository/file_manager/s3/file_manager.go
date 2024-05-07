package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileManager struct {
	client *minio.Client
}

func (fm *FileManager) Get(ctx context.Context, dir, name string) (io.ReadCloser, error) {
	return nil, app_errors.ErrUnimplemented
}

func New(cfg *config.S3) (domain.FileManager, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &FileManager{client: client}, nil
}

func (fm *FileManager) Save(ctx context.Context, bucket, name string, input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	_, err = fm.client.PutObject(ctx, bucket, name, buffer, int64(buffer.Len()), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.StatusCode == http.StatusNotFound {
			return app_errors.ErrNotFound
		}
	}

	return err
}

func (fm *FileManager) Delete(ctx context.Context, bucket, name string) error {
	err := fm.client.RemoveObject(ctx, bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.StatusCode == http.StatusNotFound {
			return app_errors.ErrNotFound
		}
	}
	return err
}
