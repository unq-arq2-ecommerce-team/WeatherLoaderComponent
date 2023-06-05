package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type SaveCurrentWeatherUseCase struct {
	baseLogger             domain.Logger
	GetCurrentWeatherQuery *GetCurrentWeatherQuery
	SaveWeatherCommand     *SaveWeatherCommand
}

func NewSaveCurrentWeatherUseCase(baseLogger domain.Logger, getCurrentWeatherQuery *GetCurrentWeatherQuery, saveWeatherCommand *SaveWeatherCommand) *SaveCurrentWeatherUseCase {
	return &SaveCurrentWeatherUseCase{
		baseLogger:             baseLogger.WithFields(domain.LoggerFields{"useCase": "saveCurrentWeatherUseCase"}),
		GetCurrentWeatherQuery: getCurrentWeatherQuery,
		SaveWeatherCommand:     saveWeatherCommand,
	}
}

func (u *SaveCurrentWeatherUseCase) Do(ctx context.Context) error {
	weather, err := u.GetCurrentWeatherQuery.Do(ctx)
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
