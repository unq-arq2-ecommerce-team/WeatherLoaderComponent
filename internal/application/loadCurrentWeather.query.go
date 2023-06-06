package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type LoadCurrentWeatherQuery struct {
	weatherRepository domain.WeatherRemoteRepository
}

func NewLoadCurrentWeatherUseCase(weatherRepository domain.WeatherRemoteRepository) *LoadCurrentWeatherQuery {
	return &LoadCurrentWeatherQuery{weatherRepository: weatherRepository}
}

func (u *LoadCurrentWeatherQuery) Do(ctx context.Context) (*domain.Weather, error) {
	return u.weatherRepository.GetCurrentWeather(ctx)
}
