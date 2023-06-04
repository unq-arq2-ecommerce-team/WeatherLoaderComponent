package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type weatherRemoteRepository struct {
	apiKey string
	apiUrl string
	lat    string
	long   string
	logger *log.Logger
}

func NewWeatherRemoteRepository(apiKey, apiUrl, lat, long string, logger *log.Logger) domain.WeatherRemoteRepository {
	return &weatherRemoteRepository{
		apiKey: apiKey,
		apiUrl: apiUrl,
		lat:    lat,
		long:   long,
		logger: logger,
	}
}

func (r *weatherRemoteRepository) GetCurrentWeather() (*domain.Weather, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&units=metric&appid=%s", r.apiUrl, r.lat, r.long, r.apiKey)
	r.logger.Println("Requesting API:", url)
	resp, err := http.Get(url)
	now := time.Now()
	if err != nil {
		r.logger.Println("Failed to request API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Println("Error reading API response:", err)
		return nil, err
	}
	r.logger.Println(string(body))

	var weatherData app.WeatherDTO
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		r.logger.Println("Error parsing API response:", err)
		return nil, err
	}
	return domain.NewWeather(weatherData.Name, weatherData.Main.Temp, weatherData.Main.FeelLike, weatherData.Main.TempMin, weatherData.Main.TempMax, weatherData.Main.Pressure, weatherData.Main.Humidity, now), nil
}
