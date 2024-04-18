package app_errors

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) error {
	var validationErrors validator.ValidationErrors
	var ok bool

	validationErrors, ok = err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errors := make(map[string][]string, 2)

	for _, err := range validationErrors {
		errors[err.Field()] = append(errors[err.Field()], err.Error())
	}

	return &ResponseError{
		Code:   http.StatusBadRequest,
		Msg:    "validation error",
		Errors: errors,
	}
}
