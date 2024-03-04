package validator

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-starter/response"
	"net/http"
)

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "positiveNumber":
		return "only positive number"
	case "domainName":
		return "domain name is not valid"

	}
	return fe.Error()
}

func HandleError(ctx *gin.Context, err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]response.ApiError, len(ve))
		for _, fe := range ve {
			out = append(out, response.ApiError{Code: 400, Message: fe.Field() + " " + msgForTag(fe)})
		}
		response.Error(http.StatusBadRequest, ctx, out, nil)
	}
	return err.Error()
}
