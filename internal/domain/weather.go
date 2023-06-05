package domain

import (
	"time"
)

type Weather struct {
	City           string    `json:"city"`
	Temperature    float64   `json:"temperature"`
	FeelsLike      float64   `json:"feelsLike"`
	TemperatureMin float64   `json:"temperatureMin"`
	TemperatureMax float64   `json:"temperatureMax"`
	Pressure       int       `json:"pressure"`
	Humidity       int       `json:"humidity"`
	TimeStamp      time.Time `json:"timestamp"`
}

func NewWeather(city string, temperature float64, feelsLike float64, temperatureMin float64, temperatureMax float64, pressure int, humidity int, timeStamp time.Time) *Weather {
	return &Weather{
		City:           city,
		Temperature:    temperature,
		FeelsLike:      feelsLike,
		TemperatureMin: temperatureMin,
		TemperatureMax: temperatureMax,
		Pressure:       pressure,
		Humidity:       humidity,
		TimeStamp:      timeStamp,
	}
}

func (w Weather) String() string {
	return ParseStruct(w)
}
