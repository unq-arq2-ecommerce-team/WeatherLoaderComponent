package dto

import (
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"time"
)

type CurrentTemperatureDTO struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp" time_format:"2006-01-02T15:04:05.000Z" time_utc:"1"`
}

func NewCurrentTemperatureDTO(currentWeather *domain.Weather) *CurrentTemperatureDTO {
	return &CurrentTemperatureDTO{
		City:        currentWeather.City,
		Temperature: currentWeather.Temperature,
		Timestamp:   currentWeather.TimeStamp,
	}
}
