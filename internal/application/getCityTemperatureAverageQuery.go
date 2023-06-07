package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"time"
)

type GetCityTemperatureAverageQuery struct {
	weatherRepository domain.WeatherLocalRepository
}

func NewGetCityTemperatureAverageQuery(weatherRepository domain.WeatherLocalRepository) *GetCityTemperatureAverageQuery {
	return &GetCityTemperatureAverageQuery{
		weatherRepository: weatherRepository,
	}
}

func (q *GetCityTemperatureAverageQuery) Do(ctx context.Context, city string, dateFrom, dateTo time.Time) (*domain.AverageTemperature, error) {
	return q.weatherRepository.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
}
