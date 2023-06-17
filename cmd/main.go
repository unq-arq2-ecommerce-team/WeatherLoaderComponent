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
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"time"
)

func main() {
	conf := config.LoadConfig()
	logger := loggerPkg.New(&loggerPkg.Config{
		ServiceName:     config.ServiceName,
		EnvironmentName: conf.Environment,
		LogLevel:        conf.LogLevel,
		LogFormat:       loggerPkg.JsonFormat,
		LokiHost:        conf.LokiHost,
	})

	mongoDb := _mongo.Connect(context.Background(), logger, conf.MongoURI, conf.MongoDatabase)

	// OTEL
	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	// domain repositories
	weatherLocalRepository := _mongo.NewWeatherLocalRepository(mongoDb, logger, conf.MongoTimeout)
	weatherRemoteRepository := http.NewWeatherRemoteRepository(logger, http.NewClient(logger, conf.Weather.HttpConfig), conf.Weather.ApiKey, conf.Weather.ApiUrl, conf.Weather.Lat, conf.Weather.Long)

	// use cases
	saveCurrentWeatherUseCase := createSaveCurrentWeatherUseCase(logger, weatherLocalRepository, weatherRemoteRepository)
	findCityCurrentTemperatureQuery := app.NewFindCityCurrentTemperatureQuery(weatherLocalRepository)
	getCityDayTemperatureAverageQuery := app.NewGetCityTemperatureAverageQuery(weatherLocalRepository)

	go startTickerOfSaveCurrentWeatherUseCase(logger, conf.TickerLoopTime, saveCurrentWeatherUseCase)

	_app := infra.NewGinApplication(conf, logger, findCityCurrentTemperatureQuery, getCityDayTemperatureAverageQuery)
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

func createSaveCurrentWeatherUseCase(logger domain.Logger, weatherLocalRepo domain.WeatherLocalRepository, weatherRemoteRepo domain.WeatherRemoteRepository) *app.SaveCurrentWeatherUseCase {
	loadCurrentWeatherQuery := app.NewLoadCurrentWeatherUseCase(weatherRemoteRepo)
	saveWeatherQuery := app.NewSaveWeatherCommand(weatherLocalRepo)
	return app.NewSaveCurrentWeatherUseCase(logger, loadCurrentWeatherQuery, saveWeatherQuery)
}

func initTracerAuto() func(context.Context) error {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
		),
	)

	if err != nil {
		log.Fatal("Could not set exporter: ", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "weather-loader"),
			attribute.String("application", "WeatherLoaderComponent"),
		),
	)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return exporter.Shutdown
}
