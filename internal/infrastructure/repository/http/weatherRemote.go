package http

import (
	"context"
	"encoding/json"
	"fmt"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
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
		baseLogger: baseLogger.WithFields(loggerPkg.Fields{"http.repository": "weatherRemoteRepository"}),
		apiKey:     apiKey,
		apiUrl:     apiUrl,
		lat:        lat,
		long:       long,
	}
}

func (r *weatherRemoteRepository) GetCurrentWeather(ctx context.Context) (*domain.Weather, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&units=metric&appid=%s", r.apiUrl, r.lat, r.long, r.apiKey)
	logger := r.baseLogger.WithFields(loggerPkg.Fields{"url": url})
	now := time.Now().UTC()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.WithFields(loggerPkg.Fields{"error": err}).Errorf("error when create request object")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		logger.WithFields(loggerPkg.Fields{"error": err}).Errorf("error when do http request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithFields(loggerPkg.Fields{"error": err}).Errorf("failed when read response body")
		return nil, err
	}

	logger = logger.WithFields(loggerPkg.Fields{"responseBody": string(body)})

	if !IsStatusCode2XX(resp.StatusCode) {
		_err := fmt.Errorf("invalid response status code %v", resp.StatusCode)
		logger.WithFields(loggerPkg.Fields{"error": _err}).Errorf("status code not in [200,299]")
		return nil, _err
	}

	var weatherData app.WeatherDTO
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		logger.WithFields(loggerPkg.Fields{"error": err}).Errorf("error parsing response body to model dto data")
		return nil, err
	}
	logger.Debugf("successful get current weather")
	return domain.NewWeather(weatherData.Name, weatherData.Main.Temp, weatherData.Main.FeelLike, weatherData.Main.TempMin, weatherData.Main.TempMax, weatherData.Main.Pressure, weatherData.Main.Humidity, now), nil
}
