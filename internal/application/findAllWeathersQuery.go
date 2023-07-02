package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type GetFindAllWeathersQuery struct {
	weatherRepository domain.WeatherLocalRepository
}

func NewGetFindAllWeathersQuery(weatherRepository domain.WeatherLocalRepository) *GetFindAllWeathersQuery {
	return &GetFindAllWeathersQuery{
		weatherRepository: weatherRepository,
	}
}

func (q *GetFindAllWeathersQuery) Do(ctx context.Context) (*[]domain.Weather, error) {
	return q.weatherRepository.FindAllWeathers(ctx)
}
