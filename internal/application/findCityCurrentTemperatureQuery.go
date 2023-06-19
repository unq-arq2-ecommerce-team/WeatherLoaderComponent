package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type FindCityCurrentTemperatureQuery struct {
	logger                    domain.Logger
	weatherLocalRepository    domain.WeatherLocalRepository
	weatherExternalRepository domain.WeatherRemoteRepository
}

func NewFindCityCurrentTemperatureQuery(logger domain.Logger, weatherLocalRepository domain.WeatherLocalRepository, weatherExternalRepository domain.WeatherRemoteRepository) *FindCityCurrentTemperatureQuery {
	return &FindCityCurrentTemperatureQuery{
		logger:                    logger.WithFields(domain.LoggerFields{"action.query": "FindCityCurrentTemperatureQuery"}),
		weatherLocalRepository:    weatherLocalRepository,
		weatherExternalRepository: weatherExternalRepository,
	}
}

func (q *FindCityCurrentTemperatureQuery) Do(ctx context.Context, city string) (*domain.Weather, error) {
	log := q.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"city": city})
	weather, err := q.weatherLocalRepository.FindCurrentByCity(ctx, city)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when FindCurrentByCity %s from weatherLocalRepository", city)
		weather, err = q.doFallback(ctx, log, city)
	}
	log.Info("successful get current temperature weather with city %s", city)
	return weather, err
}

func (q *FindCityCurrentTemperatureQuery) doFallback(ctx context.Context, log domain.Logger, city string) (*domain.Weather, error) {
	log.Debugf("init fallback")
	lat, long, err := q.weatherExternalRepository.GetLatAndLongFromCity(ctx, city)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error fallback when GetLatAndLongFromCity %s from weatherExternalRepository", city)
		return nil, err
	}
	log.Debug("successful get current temperature weather with city %s from fallback", city)
	return q.weatherExternalRepository.GetCurrentWeather(ctx, lat, long)
}
