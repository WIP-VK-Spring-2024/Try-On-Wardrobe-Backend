package api_errors

import "errors"

var (
	ErrNotFound              = errors.New("requested resource does not exist")
	ErrInvalidToken          = errors.New("invalid token format")
	ErrExpired               = errors.New("token expired")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAlreadyExists         = errors.New("resource already exists")
	ErrSessionNotInitialized = errors.New("failed initializing session")
)
