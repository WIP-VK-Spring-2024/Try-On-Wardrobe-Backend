package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var usernameRegexp = regexp.MustCompile(`^[a-zA-Zа-яёА-ЯЁ0-9-_()+=~@^:?;$#№%*@|{}[\]!<>]+$`)

func init() {
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegexp.MatchString(fl.Field().String())
	})
}

func Struct(item any) error {
	return validate.Struct(item)
}
