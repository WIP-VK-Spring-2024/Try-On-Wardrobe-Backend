package api_errors

import "errors"

var (
	ErrNotFound     = errors.New("requested resource does not exist")
	ErrInvalidToken = errors.New("invalid token format")
	ErrExpired      = errors.New("token expired")
)
