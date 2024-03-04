package validator

import "github.com/go-playground/validator/v10"

var positiveNumber validator.Func = func(fl validator.FieldLevel) bool {
	integer, ok := fl.Field().Interface().(int32)
	if ok {
		if integer <= 0 {
			return false
		}
	}
	return true
}
