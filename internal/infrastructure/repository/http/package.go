package http

import (
	"context"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/otel"
	"net/http"
	"time"
)

const deltaRetryWait = 2 * time.Second

func NewDefaultClient() *http.Client {
	return cleanhttp.DefaultPooledClient()
}

func NewClient(logger domain.Logger, httpConfig config.HttpConfig) *http.Client {
	httpClient := NewDefaultClient()
	httpClient.Timeout = httpConfig.Timeout
	if httpConfig.OtelEnabled {
		httpClient.Transport = otel.WrapAndReturn(httpClient.Transport)
	}

	retryableClient := retryablehttp.NewClient()
	retryableClient.Logger = logger.WithFields(domain.LoggerFields{"loggerFrom": "http.retryableClient"})
	retryableClient.HTTPClient = httpClient
	retryableClient.RetryMax = httpConfig.Retries
	retryableClient.RetryWaitMin = httpConfig.RetryWait
	retryableClient.RetryWaitMax = httpConfig.RetryWait + deltaRetryWait
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
