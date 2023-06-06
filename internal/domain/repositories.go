package domain

import "context"

type WeatherLocalRepository interface {
	Save(context.Context, *Weather) error
	FindCurrentByCity(ctx context.Context, city string) (*Weather, error)
}

type WeatherRemoteRepository interface {
	GetCurrentWeather(context.Context) (*Weather, error)
}
