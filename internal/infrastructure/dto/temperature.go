package dto

import (
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"time"
)

type CurrentTemperatureDTO struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewCurrentTemperatureDTO(currentWeather *domain.Weather) *CurrentTemperatureDTO {
	return &CurrentTemperatureDTO{
		City:        currentWeather.City,
		Temperature: currentWeather.Temperature,
		Timestamp:   currentWeather.TimeStamp,
	}
}
