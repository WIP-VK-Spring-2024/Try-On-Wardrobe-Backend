package utils

import (
	"try-on/internal/pkg/api_errors"

	"gorm.io/gorm"
)

func GormError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return api_errors.ErrNotFound
	case gorm.ErrDuplicatedKey:
		return api_errors.ErrAlreadyExists
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
