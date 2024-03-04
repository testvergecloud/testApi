package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators all custom validators must be registered here
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positiveNumber", positiveNumber)
		v.RegisterValidation("domainName", domainName)
	}
}
