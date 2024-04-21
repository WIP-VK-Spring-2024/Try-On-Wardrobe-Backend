package app_errors

import (
	"errors"
	"net/http"
	"runtime"
)

var (
	ErrNotFound                = errors.New("requested resource does not exist")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrAlreadyExists           = errors.New("resource already exists")
	ErrTokenMalformed          = errors.New("token malformed or missing")
	ErrInvalidSignature        = errors.New("token has invalid signature")
	ErrTokenExpired            = errors.New("token has expired")
	ErrUnimplemented           = errors.New("method unimplemented")
	ErrNoRelatedEntity         = errors.New("related resource not found")
	ErrConstraintViolation     = errors.New("constraint violated")
	ErrNotOwner                = errors.New("must be the owner to delete or edit this resource")
	ErrTryOnInvalidClothesNum  = errors.New("try on requires at least 1 garment, but not more than 2")
	ErrTryOnInvalidClothesType = errors.New("try on requires 1 dress, or a maximum of 1 of upper body and 1 lower body garments")
	ErrNotEnoughClothes        = errors.New("outfit generation requires at least 1 upper and 1 lower garment")
)

var (
	ErrBadRequest = &ResponseError{
		Msg:  "bad request",
		Code: http.StatusBadRequest,
	}

	ErrUnauthorized = &ResponseError{
		Msg:  "credentials missing or invalid",
		Code: http.StatusUnauthorized,
	}

	ErrUserIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "user ID is missing or isn't a valid uuid",
	}

	ErrClothesIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "clothes ID is missing or isn't a valid uuid",
	}

	ErrUserImageIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "user image ID is missing or isn't a valid uuid",
	}

	ErrTryOnIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "try on result ID is missing or isn't a valid uuid",
	}

	ErrOutfitIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "outfit ID is missing or isn't a valid uuid",
	}

	ErrPostIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "post ID is missing or isn't a valid uuid",
	}

	ErrCommentIdInvalid = &ResponseError{
		Code: http.StatusBadRequest,
		Msg:  "comment ID is missing or isn't a valid uuid",
	}
)

type InternalError struct {
	Err  error
	File string
	Line int
}

func (err *InternalError) Error() string {
	return err.Err.Error()
}

//easyjson:json
type ResponseError struct {
	Code   int `json:"-"`
	Msg    string
	Errors map[string][]string
}

func (err ResponseError) Error() string {
	return err.Msg
}

func New(err error) error {
	var code int

	switch {
	case errors.Is(err, ErrAlreadyExists):
		code = http.StatusConflict

	case errors.Is(err, ErrNotOwner) || errors.Is(err, ErrInvalidCredentials):
		code = http.StatusForbidden

	case errors.Is(err, ErrNotFound) || errors.Is(err, ErrNoRelatedEntity):
		code = http.StatusNotFound

	case Any(err, ErrConstraintViolation, ErrTryOnInvalidClothesNum,
		ErrTryOnInvalidClothesType, ErrNotEnoughClothes):
		code = http.StatusBadRequest

	default:
		_, file, line, _ := runtime.Caller(1)
		return &InternalError{
			Err:  err,
			File: file,
			Line: line,
		}
	}

	return &ResponseError{
		Code: code,
		Msg:  err.Error(),
	}
}

func Any(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
