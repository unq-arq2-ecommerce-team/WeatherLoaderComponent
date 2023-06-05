package domain

import "context"

type WeatherLocalRepository interface {
	Save(context.Context, *Weather) error
}

type WeatherRemoteRepository interface {
	GetCurrentWeather(context.Context) (*Weather, error)
}
