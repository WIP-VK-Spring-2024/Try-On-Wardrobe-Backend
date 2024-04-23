package validate

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var usernameRegexp = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9-_()+=~@^:?;$#№%*@|{}[\]!<>]+$`)

func init() {
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		fmt.Printf("Validating value: %q", fl.Field().String())
		return usernameRegexp.MatchString(fl.Field().String())
	})
}

func Struct(item any) error {
	return validate.Struct(item)
}
