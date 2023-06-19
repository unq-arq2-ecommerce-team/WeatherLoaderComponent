package domain

import (
	"context"
	"time"
)

type WeatherLocalRepository interface {
	Save(context.Context, *Weather) error
	FindCurrentByCity(ctx context.Context, city string) (*Weather, error)
	GetAverageTemperatureByCityAndDateRange(ctx context.Context, city string, dateFrom, dateTo time.Time) (*AverageTemperature, error)
}

type WeatherRemoteRepository interface {
	GetCurrentWeather(ctx context.Context, lat, long string) (*Weather, error)
	GetLatAndLongFromCity(ctx context.Context, city string) (lat, long string, err error)
}
