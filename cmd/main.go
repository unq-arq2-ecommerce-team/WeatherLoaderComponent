package main

import (
	"context"
	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	infra "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/repository/http"
	_mongo "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/repository/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func main() {
	conf := config.LoadConfig()
	logger := loggerPkg.New(&loggerPkg.Config{
		ServiceName:     config.ServiceName,
		EnvironmentName: conf.Environment,
		LogLevel:        conf.LogLevel,
		LogFormat:       loggerPkg.JsonFormat,
	})

	mongoDb := _mongo.Connect(context.Background(), logger, conf.MongoURI, conf.MongoDatabase)
	saveCurrentWeatherUseCase := createSaveCurrentWeatherUseCase(mongoDb, conf.MongoTimeout, conf.Weather, logger)

	go startTickerOfSaveCurrentWeatherUseCase(logger, conf.TickerLoopTime, saveCurrentWeatherUseCase)

	_app := infra.NewGinApplication(conf, logger)
	logger.Fatal(_app.Run())
}

// startTickerOfSaveCurrentWeatherUseCase init a job which runs periodically use case param every duration of tickerLoopTime param
func startTickerOfSaveCurrentWeatherUseCase(baseLogger domain.Logger, tickerLoopTime time.Duration, useCase *app.SaveCurrentWeatherUseCase) {
	logger := baseLogger.WithFields(domain.LoggerFields{"logger": "ticker"})
	ticker := time.NewTicker(tickerLoopTime)
	logger.Infof("starting ticker loop with loop duration %s", tickerLoopTime.String())

	useCaseDoAndLogErrFn := func() {
		if err := useCase.Do(context.Background()); err != nil {
			logger.WithFields(domain.LoggerFields{"error": err}).Errorf("Error saving current weather")
		}
	}

	// Start off by calling API immediately.
	useCaseDoAndLogErrFn()
	for range ticker.C {
		useCaseDoAndLogErrFn()
	}
}

func createSaveCurrentWeatherUseCase(mongoConn *mongo.Database, mongoTimeout time.Duration, weatherConf config.Weather, logger domain.Logger) *app.SaveCurrentWeatherUseCase {
	localRepository := _mongo.NewWeatherLocalRepository(mongoConn, logger, mongoTimeout)
	remoteRepository := http.NewWeatherRemoteRepository(logger, http.NewClient(), weatherConf.ApiKey, weatherConf.ApiUrl, weatherConf.Lat, weatherConf.Long)
	getCurrentWeatherQuery := app.NewGetCurrentWeatherUseCase(remoteRepository)
	saveWeatherQuery := app.NewSaveWeatherCommand(localRepository)
	return app.NewSaveCurrentWeatherUseCase(logger, getCurrentWeatherQuery, saveWeatherQuery)
}
