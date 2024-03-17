package utils

import (
	"errors"

	"try-on/internal/pkg/app_errors"

	"gorm.io/gorm"
)

func GormError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return errors.Join(err, app_errors.ErrNotFound)
	case gorm.ErrDuplicatedKey:
		return errors.Join(err, app_errors.ErrAlreadyExists)
	case gorm.ErrInvalidField:
		return errors.Join(err, app_errors.ErrBadRequest)
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
