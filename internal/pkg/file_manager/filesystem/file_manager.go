package filesystem

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
)

type FileManager struct {
	baseDir string
}

func New(baseDir string) domain.FileManager {
	return &FileManager{
		baseDir: baseDir,
	}
}

func (fm *FileManager) Save(_ context.Context, dir, name string, input io.Reader) error {
	err := os.MkdirAll(fm.baseDir+"/"+dir, fs.ModePerm)
	if err != nil {
		return err
	}

	newFile, err := os.Create(fm.fullName(dir, name))
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, input)
	if err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) Delete(_ context.Context, dir, name string) error {
	err := os.Remove(fm.fullName(dir, name))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		err = app_errors.ErrNotFound
	}
	return err
}

func (fm *FileManager) fullName(dir, name string) string {
	filepath := fm.baseDir + "/"
	if dir != "" {
		filepath += dir + "/"
	}
	return filepath + name
}
