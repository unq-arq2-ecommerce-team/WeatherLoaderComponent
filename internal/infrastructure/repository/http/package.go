package http

import (
	"github.com/hashicorp/go-cleanhttp"
	"net/http"
)

func NewClient() *http.Client {
	return cleanhttp.DefaultPooledClient()
}

func IsStatusCode2XX(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
