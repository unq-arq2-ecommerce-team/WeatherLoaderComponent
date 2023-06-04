package domain

type WeatherLocalRepository interface {
	Save(weather *Weather) error
}

type WeatherRemoteRepository interface {
	GetCurrentWeather() (*Weather, error)
}
