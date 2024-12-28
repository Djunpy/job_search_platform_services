package server

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CreateAndSendRequest(method, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	client := &http.Client{}
	return client.Do(req)
}

func GetReqFullUrl(ctx *gin.Context, target string) string {
	fullPath := ctx.Request.URL.Path
	rawQuery := ctx.Request.URL.RawQuery
	url := target + fullPath
	if rawQuery != "" {
		url += "?" + rawQuery
	}
	return url
}
