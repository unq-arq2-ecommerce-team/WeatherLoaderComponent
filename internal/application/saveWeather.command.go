package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type SaveWeatherCommand struct {
	weatherRepository domain.WeatherLocalRepository
}

func NewSaveWeatherCommand(weatherRepository domain.WeatherLocalRepository) *SaveWeatherCommand {
	return &SaveWeatherCommand{weatherRepository: weatherRepository}
}

func (u *SaveWeatherCommand) Do(ctx context.Context, weather *domain.Weather) error {
	return u.weatherRepository.Save(ctx, weather)
}
