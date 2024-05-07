package validate

import (
	"regexp"

	"try-on/internal/pkg/utils"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var (
	usernameRegexp   = regexp.MustCompile(`^[a-zA-Zа-яёА-ЯЁ0-9-_()+=~@^:?;$#№%*@|{}[\]!<>]+$`)
	otherNamesRegexp = regexp.MustCompile(`^[a-zA-Zа-яёА-ЯЁ0-9 -_()+=~@^:?;$#№%*@|{}[\]!<>]+$`)
)

func init() {
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegexp.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		return otherNamesRegexp.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("name_slice", func(fl validator.FieldLevel) bool {
		slice, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}
		return utils.Every(slice, otherNamesRegexp.MatchString)
	})
}

func Struct(item any) error {
	return validate.Struct(item)
}
