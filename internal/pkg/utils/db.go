package utils

import (
	"errors"

	"try-on/internal/pkg/app_errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func PgxError(err error) error {
	if err == pgx.ErrNoRows {
		return app_errors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	var appError error

	switch {
	case pgErr.Code == pgerrcode.UniqueViolation:
		appError = app_errors.ErrAlreadyExists
	case pgErr.Code == pgerrcode.NoData:
		appError = app_errors.ErrNotFound
	case pgErr.Code == pgerrcode.ForeignKeyViolation:
		appError = app_errors.ErrNoRelatedEntity
	case pgerrcode.IsIntegrityConstraintViolation(pgErr.Code):
		appError = app_errors.ErrConstraintViolation
	default:
		return err
	}

	return errors.Join(err, appError)
}
