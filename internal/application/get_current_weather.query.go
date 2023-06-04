package application

import (
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type GetCurrentWeatherQuery struct {
	weatherRepository domain.WeatherRemoteRepository
}

func NewGetCurrentWeatherUseCase(weatherRepository domain.WeatherRemoteRepository) *GetCurrentWeatherQuery {
	return &GetCurrentWeatherQuery{weatherRepository: weatherRepository}
}

func (u *GetCurrentWeatherQuery) Do() (*domain.Weather, error) {
	return u.weatherRepository.GetCurrentWeather()
}
