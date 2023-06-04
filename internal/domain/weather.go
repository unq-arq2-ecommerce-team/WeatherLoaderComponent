package domain

import "time"

type Weather struct {
	City           string
	Temperature    float64
	FeelsLike      float64
	TemperatureMin float64
	TemperatureMax float64
	Pressure       int
	Humidity       int
	TimeStamp      time.Time
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
