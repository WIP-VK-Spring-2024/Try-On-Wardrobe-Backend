package utils

import (
	"try-on/internal/pkg/app_errors"

	"gorm.io/gorm"
)

func GormError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return app_errors.ErrNotFound
	case gorm.ErrDuplicatedKey:
		return app_errors.ErrAlreadyExists
	default:
		return err
	}
}

func TranslateGormError[T any](item *T, err error) (*T, error) {
	err = GormError(err)
	if err != nil {
		return nil, err
	}
	return item, nil
}
