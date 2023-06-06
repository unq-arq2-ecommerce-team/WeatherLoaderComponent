package domain

import (
	"time"
)

type Weather struct {
	City           string    `json:"city" bson:"city"`
	Temperature    float64   `json:"temperature" bson:"temperature"`
	FeelsLike      float64   `json:"feelsLike" bson:"feelsLike"`
	TemperatureMin float64   `json:"temperatureMin" bson:"temperatureMin"`
	TemperatureMax float64   `json:"temperatureMax" bson:"temperatureMax"`
	Pressure       int       `json:"pressure" bson:"pressure"`
	Humidity       int       `json:"humidity" bson:"humidity"`
	TimeStamp      time.Time `json:"timestamp" bson:"timestamp"`
}

func NewWeather(city string, temperature, feelsLike, temperatureMin, temperatureMax float64, pressure, humidity int, timeStamp time.Time) *Weather {
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
