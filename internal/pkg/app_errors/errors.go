package app_errors

import (
	"errors"
	"net/http"
	"runtime"
)

var (
	ErrNotFound           = errors.New("requested resource does not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrTokenMalformed     = errors.New("token malformed or missing")
	ErrInvalidSignature   = errors.New("token has invalid signature")
	ErrTokenExpired       = errors.New("token has expired")
	ErrUnimplemented      = errors.New("method unimplemented")
)

var (
	ErrBadRequest = ErrorMsg{
		Msg:  "bad request",
		Code: http.StatusBadRequest,
	}
	ErrUnauthorized = ErrorMsg{
		Msg:  "credentials missing or invalid",
		Code: http.StatusUnauthorized,
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
type ErrorMsg struct {
	Code int `json:"-"`
	Msg  string
}

func (err ErrorMsg) Error() string {
	return err.Msg
}

func New(err error) error {
	code := http.StatusInternalServerError

	switch {
	default:
		_, file, line, _ := runtime.Caller(1)
		return &InternalError{
			Err:  err,
			File: file,
			Line: line,
		}
	case errors.Is(err, ErrAlreadyExists):
		code = http.StatusConflict

	case errors.Is(err, ErrAlreadyExists):
		code = http.StatusNotFound

	case errors.Is(err, ErrInvalidCredentials):
		code = http.StatusForbidden
	}

	return ErrorMsg{
		Code: code,
		Msg:  err.Error(),
	}
}
