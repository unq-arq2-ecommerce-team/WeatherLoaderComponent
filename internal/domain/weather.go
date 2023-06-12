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
	TimeStamp      time.Time `json:"timestamp" bson:"timestamp" time_format:"2006-01-02T15:04:05.000Z" time_utc:"1"`
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

type AverageTemperature struct {
	City           string    `json:"city" bson:"_id"`
	AvgTemperature float64   `json:"avgTemperature" bson:"avgTemperature"`
	DateFrom       time.Time `json:"dateFrom" time_format:"2006-01-02T15:04:05.000Z" time_utc:"1"`
	DateTo         time.Time `json:"dateTo" time_format:"2006-01-02T15:04:05.000Z" time_utc:"1"`
	DaysBetween    float64   `json:"daysBetween"`
}

func (a *AverageTemperature) Set(dateFrom, dateTo time.Time) {
	a.DateFrom = dateFrom
	a.DateTo = dateTo
	a.DaysBetween = dateTo.Sub(dateFrom).Hours() / 24
}

func (w Weather) String() string {
	return ParseStruct(w)
}

func (a *AverageTemperature) String() string {
	return ParseStruct(a)
}
