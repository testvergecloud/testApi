package response

import (
	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Errors  []ApiError  `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(code int, ctx *gin.Context, data interface{}) {
	response := ApiResponse{true, nil, data}
	ctx.JSON(code, response)
}

func Error(code int, ctx *gin.Context, errors []ApiError, data interface{}) {
	response := ApiResponse{false, errors, data}
	ctx.JSON(code, response)
}
