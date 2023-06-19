package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type SaveCurrentWeatherUseCase struct {
	baseLogger              domain.Logger
	LoadCurrentWeatherQuery *LoadCurrentWeatherQuery
	SaveWeatherCommand      *SaveWeatherCommand
}

func NewSaveCurrentWeatherUseCase(baseLogger domain.Logger, loadCurrentWeatherQuery *LoadCurrentWeatherQuery, saveWeatherCommand *SaveWeatherCommand) *SaveCurrentWeatherUseCase {
	return &SaveCurrentWeatherUseCase{
		baseLogger:              baseLogger.WithFields(domain.LoggerFields{"useCase": "saveCurrentWeatherUseCase"}),
		LoadCurrentWeatherQuery: loadCurrentWeatherQuery,
		SaveWeatherCommand:      saveWeatherCommand,
	}
}

func (u *SaveCurrentWeatherUseCase) Do(ctx context.Context, lat, long string) error {
	weather, err := u.LoadCurrentWeatherQuery.Do(ctx, lat, long)
	logger := u.baseLogger.WithFields(domain.LoggerFields{"weather": weather})
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Error("error when get current weather")
		return err
	}
	if err := u.SaveWeatherCommand.Do(ctx, weather); err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Error("error when save weather")
		return err
	}
	logger.Info("successful save current weather use case")
	return nil
}
