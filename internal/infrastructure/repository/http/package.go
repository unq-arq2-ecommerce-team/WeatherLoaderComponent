package http

import (
	"context"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"net/http"
	"time"
)

func NewDefaultClient() *http.Client {
	return cleanhttp.DefaultPooledClient()
}

func NewClient(httpConfig config.HttpConfig) *http.Client {
	httpClient := NewDefaultClient()
	httpClient.Timeout = httpConfig.Timeout

	retryableClient := retryablehttp.NewClient()
	retryableClient.HTTPClient = httpClient
	retryableClient.RetryMax = httpConfig.Retries
	retryableClient.RetryWaitMin = httpConfig.RetryWait
	retryableClient.RetryWaitMax = httpConfig.RetryWait + (2 * time.Second)
	return retryableClient.StandardClient()
}

func IsStatusCode2XX(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func NewRequestWithContextWithNoBody(ctx context.Context, httpMethod, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil)
	req.Header.Add(logger.RequestIdHeaderKey(), logger.GetRequestId(ctx))
	return req, err
}
