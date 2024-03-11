package app_errors

import (
	"errors"
	"runtime"
)

var (
	ErrNotFound              = errors.New("requested resource does not exist")
	ErrInvalidToken          = errors.New("invalid token format")
	ErrExpired               = errors.New("token expired")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAlreadyExists         = errors.New("resource already exists")
	ErrSessionNotInitialized = errors.New("failed initializing session")
	ErrTokenMalformed        = errors.New("token malformed or missing")
	ErrInvalidSignature      = errors.New("token has invalid signature")
	ErrTokenExpired          = errors.New("token has expired")
)

type Error struct {
	Err  error
	File string
	Line int
}

func (err *Error) Error() string {
	return err.Err.Error()
}

func New(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.Join(&Error{
		Err:  err,
		File: file,
		Line: line,
	}, err)
}

//easyjson:json
type ErrorMsg struct {
	Msg string
}
