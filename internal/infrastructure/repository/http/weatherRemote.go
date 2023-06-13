package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type weatherRemoteRepository struct {
	baseLogger domain.Logger
	client     *http.Client
	apiKey     string
	apiUrl     string
	lat        string
	long       string
}

func NewWeatherRemoteRepository(baseLogger domain.Logger, client *http.Client, apiKey, apiUrl, lat, long string) domain.WeatherRemoteRepository {
	return &weatherRemoteRepository{
		client:     client,
		baseLogger: baseLogger.WithFields(domain.LoggerFields{"http.repository": "weatherRemoteRepository"}),
		apiKey:     apiKey,
		apiUrl:     apiUrl,
		lat:        lat,
		long:       long,
	}
}

func (r *weatherRemoteRepository) GetCurrentWeather(ctx context.Context) (*domain.Weather, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&units=metric&appid=%s", r.apiUrl, r.lat, r.long, r.apiKey)
	logger := r.baseLogger.WithFields(domain.LoggerFields{"url": url})

	req, err := NewRequestWithContextWithNoBody(ctx, http.MethodGet, url)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when create request object")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when do http request")
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now().UTC()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("failed when read response body")
		return nil, err
	}

	logger = logger.WithFields(domain.LoggerFields{"responseBody": string(body)})

	if !IsStatusCode2XX(resp.StatusCode) {
		_err := fmt.Errorf("invalid response status code %v", resp.StatusCode)
		logger.WithFields(domain.LoggerFields{"error": _err}).Errorf("status code not in [200,299]")
		return nil, _err
	}

	var weatherRes app.WeatherDTO
	err = json.Unmarshal(body, &weatherRes)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error parsing response body to model dto data")
		return nil, err
	}
	if weatherRes.IsInvalid() {
		logger.WithFields(domain.LoggerFields{"weatherResDto": weatherRes}).Errorf("error invalid weather data dto")
		return nil, fmt.Errorf("error invalid weather response dto")
	}
	logger.Debugf("successful get current weather")
	return domain.NewWeather(weatherRes.City, weatherRes.Data.Temp, weatherRes.Data.FeelLike, weatherRes.Data.TempMin, weatherRes.Data.TempMax, weatherRes.Data.Pressure, weatherRes.Data.Humidity, now), nil
}
