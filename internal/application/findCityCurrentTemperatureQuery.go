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
		return q.getCurrentWeatherFromExternal(ctx, log, city)
	}
	log.Infof("successful get current temperature weather with city %s", city)
	return weather, err
}

// getCurrentWeatherFromExternal is a fallback
func (q *FindCityCurrentTemperatureQuery) getCurrentWeatherFromExternal(ctx context.Context, log domain.Logger, city string) (*domain.Weather, error) {
	log.Warnf("executing getCurrentWeatherFromExternal with city %s fallback", city)
	lat, long, err := q.weatherExternalRepository.GetLatAndLongFromCity(ctx, city)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error fallback when GetLatAndLongFromCity %s from weatherExternalRepository", city)
		return nil, err
	}
	log.Infof("successful get current temperature weather with city %s from fallback", city)
	return q.weatherExternalRepository.GetCurrentWeather(ctx, lat, long)
}
