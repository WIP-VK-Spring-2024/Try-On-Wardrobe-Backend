package httpuploader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
)

type FileManager struct {
	cfg *config.HttpApi
}

func New(cfg *config.HttpApi) domain.FileManager {
	return &FileManager{
		cfg: cfg,
	}
}

func httpOk(code int) bool {
	return code >= 200 && code < 300
}

func (fm *FileManager) Save(ctx context.Context, dir, name string, input io.Reader) error {
	payload := &bytes.Buffer{}
	form := multipart.NewWriter(payload)

	file, err := form.CreateFormFile("file", dir+"/"+name)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, input)
	if err != nil {
		return err
	}

	form.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fm.cfg.Endpoint+fm.cfg.UploadUrl, payload)
	if err != nil {
		return err
	}

	req.Header.Set(fm.cfg.TokenHeader, fm.cfg.Token)
	req.Header.Set(common.HeaderContentType, form.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !httpOk(resp.StatusCode) {
		return errors.New(resp.Status)
	}
	return err
}

func (fm *FileManager) Delete(ctx context.Context, dir, name string) error {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodDelete,
		fm.cfg.Endpoint+fm.cfg.DeleteUrl+"/"+name+"?folder="+dir,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set(fm.cfg.TokenHeader, fm.cfg.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !httpOk(resp.StatusCode) {
		return errors.New(resp.Status)
	}
	return err
}
