package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	"io"
	"net/http"
	"strings"
	"time"

	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type weatherRemoteRepository struct {
	baseLogger   domain.Logger
	client       *http.Client
	apiKey       string
	apiUrl       string
	geocodingUrl string
}

func NewWeatherRemoteRepository(baseLogger domain.Logger, client *http.Client, conf config.Weather) domain.WeatherRemoteRepository {
	return &weatherRemoteRepository{
		client:       client,
		baseLogger:   baseLogger.WithFields(domain.LoggerFields{"http.repository": "weatherRemoteRepository"}),
		apiKey:       conf.ApiKey,
		apiUrl:       conf.ApiUrl,
		geocodingUrl: conf.GeocodingUrl,
	}
}

func (r *weatherRemoteRepository) GetCurrentWeather(ctx context.Context, lat, long string) (*domain.Weather, error) {
	url := r.replaceApiKeyFromUrl(r.apiUrl)
	url = strings.Replace(url, "{LAT}", lat, -1)
	url = strings.Replace(url, "{LONG}", long, -1)
	logger := r.baseLogger.WithRequestId(ctx).WithFields(domain.LoggerFields{"url": url, "method": "GetCurrentWeather"})

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

func (r *weatherRemoteRepository) GetLatAndLongFromCity(ctx context.Context, city string) (lat, long string, err error) {
	url := r.replaceApiKeyFromUrl(r.geocodingUrl)
	url = strings.Replace(url, "{city}", city, -1)
	logger := r.baseLogger.WithRequestId(ctx).WithFields(domain.LoggerFields{"url": url, "method": "GetLatAndLongFromCity"})

	lat, long, err = "", "", nil
	req, err := NewRequestWithContextWithNoBody(ctx, http.MethodGet, url)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when create request object")
		return
	}

	resp, err := r.client.Do(req)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when do http request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("failed when read response body")
		return
	}

	logger = logger.WithFields(domain.LoggerFields{"responseBody": string(body)})

	if !IsStatusCode2XX(resp.StatusCode) {
		err = fmt.Errorf("invalid response status code %v", resp.StatusCode)
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("status code not in [200,299]")
		return
	}

	var geocodingRes []app.GeoCodingDTO
	err = json.Unmarshal(body, &geocodingRes)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error parsing response body to model dto data")
		return
	}
	if len(geocodingRes) == 0 {
		err = fmt.Errorf("geocoding response is empty")
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error parsing response body to model dto data")
		return
	}
	lat, long = fmt.Sprintf("%v", geocodingRes[0].Lat), fmt.Sprintf("%v", geocodingRes[0].Lon)
	logger.Infof("successful get geocoding from city %s with lat %s and long %s", city, lat, long)
	return
}

func (r *weatherRemoteRepository) replaceApiKeyFromUrl(url string) string {
	return strings.Replace(url, "{API_KEY}", r.apiKey, -1)
}
