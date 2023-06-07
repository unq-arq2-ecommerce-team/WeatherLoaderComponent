package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type FindCityCurrentTemperatureQuery struct {
	weatherRepository domain.WeatherLocalRepository
}

func NewFindCityCurrentTemperatureQuery(weatherRepository domain.WeatherLocalRepository) *FindCityCurrentTemperatureQuery {
	return &FindCityCurrentTemperatureQuery{
		weatherRepository: weatherRepository,
	}
}

func (q *FindCityCurrentTemperatureQuery) Do(ctx context.Context, city string) (*domain.Weather, error) {
	return q.weatherRepository.FindCurrentByCity(ctx, city)
}
