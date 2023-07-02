package main

import (
	"context"
	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	infra "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/otel"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/repository/http"
	_mongo "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/repository/mongo"
	"time"
)

func main() {
	conf := config.LoadConfig()
	isIntegrationEnv := conf.IsIntegrationEnv()

	logger := loggerPkg.New(&loggerPkg.Config{
		ServiceName:      config.ServiceName,
		EnvironmentName:  conf.Environment,
		IsIntegrationEnv: isIntegrationEnv,
		LogLevel:         conf.LogLevel,
		LogFormat:        loggerPkg.JsonFormat,
	})

	mongoDb := _mongo.Connect(context.Background(), logger, conf.Mongo.URI, conf.Mongo.Database, isIntegrationEnv)

	// OTEL
	cleanupFn := otel.InitOtelTrace(context.Background(), logger, conf.Otel, isIntegrationEnv)
	defer cleanupFn()

	// domain repositories
	conf.Weather.HttpConfig.OtelEnabled = isIntegrationEnv
	weatherLocalRepository := _mongo.NewWeatherLocalRepository(mongoDb, logger, conf.Mongo.Timeout)
	weatherRemoteRepository := http.NewWeatherRemoteRepository(logger, http.NewClient(logger, conf.Weather.HttpConfig), conf.Weather)

	// use cases
	saveCurrentWeatherUseCase := createSaveCurrentWeatherUseCase(logger, weatherLocalRepository, weatherRemoteRepository)
	findCityCurrentTemperatureQuery := app.NewFindCityCurrentTemperatureQuery(logger, weatherLocalRepository, weatherRemoteRepository)
	getCityDayTemperatureAverageQuery := app.NewGetCityTemperatureAverageQuery(weatherLocalRepository)
	getFindAllWeathersQuery := app.NewGetFindAllWeathersQuery(weatherLocalRepository)

	go startTickerOfSaveCurrentWeatherUseCase(logger, conf.TickerLoopTime, saveCurrentWeatherUseCase, conf.Weather.Lat, conf.Weather.Long)

	_app := infra.NewGinApplication(conf, logger, mongoDb.Client(), findCityCurrentTemperatureQuery, getCityDayTemperatureAverageQuery, getFindAllWeathersQuery)
	logger.Fatal(_app.Run())
}

// startTickerOfSaveCurrentWeatherUseCase init a job which runs periodically use case param every duration of tickerLoopTime param
func startTickerOfSaveCurrentWeatherUseCase(baseLogger domain.Logger, tickerLoopTime time.Duration, useCase *app.SaveCurrentWeatherUseCase, lat, long string) {
	logger := baseLogger.WithFields(domain.LoggerFields{"logger": "ticker"})
	ticker := time.NewTicker(tickerLoopTime)
	logger.Infof("starting ticker loop with loop duration %s", tickerLoopTime.String())

	useCaseDoAndLogErrFn := func() {
		if err := useCase.Do(context.Background(), lat, long); err != nil {
			logger.WithFields(domain.LoggerFields{"error": err}).Errorf("Error saving current weather")
		}
	}

	// Start off by calling API immediately.
	useCaseDoAndLogErrFn()
	for range ticker.C {
		useCaseDoAndLogErrFn()
	}
}

func createSaveCurrentWeatherUseCase(logger domain.Logger, weatherLocalRepo domain.WeatherLocalRepository, weatherRemoteRepo domain.WeatherRemoteRepository) *app.SaveCurrentWeatherUseCase {
	loadCurrentWeatherQuery := app.NewLoadCurrentWeatherUseCase(weatherRemoteRepo)
	saveWeatherQuery := app.NewSaveWeatherCommand(weatherLocalRepo)
	return app.NewSaveCurrentWeatherUseCase(logger, loadCurrentWeatherQuery, saveWeatherQuery)
}
