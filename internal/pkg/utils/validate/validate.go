package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	usernameRegexp := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9-_()+=~@^:?;$#№%*@|{}[\]!<>]+$`)

	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegexp.MatchString(fl.Field().String())
	})
}

func Struct(item any) error {
	return validate.Struct(item)
}
