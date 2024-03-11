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
)

type Error struct {
	Err  error
	File string
	Line int
}

func (err *Error) Error() string {
	return err.Err.Error()
}

func New(err error) *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Err:  err,
		File: file,
		Line: line,
	}
}

//easyjson:json
type ErrorMsg struct {
	Msg string
}
