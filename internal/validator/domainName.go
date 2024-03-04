package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var domainName validator.Func = func(fl validator.FieldLevel) bool {
	domainName := fl.Field().String()
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9-\.]+$`)
	return domainRegex.MatchString(domainName)
}
