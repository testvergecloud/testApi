package utils

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func SendRequest(ctx *gin.Context, target string, authHeader string) {

	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	rawQuery := ctx.Request.URL.RawQuery

	url := target + path + "?" + rawQuery

	proxyReq, err := http.NewRequest(method, url, ctx.Request.Body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	proxyReq.Header.Add("Authorization", authHeader)

	for name, values := range ctx.Request.Header {
		if name == "Accept-Encoding" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		// todo: handle errors depends on its type
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	result := Sanitize(body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), result)
}
