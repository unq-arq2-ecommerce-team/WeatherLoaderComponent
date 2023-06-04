package application

import "log"

type SaveCurrentWeatherUseCase struct {
	GetCurrentWeatherQuery *GetCurrentWeatherQuery
	SaveWeatherCommand     *SaveWeatherCommand
}

func NewSaveCurrentWeatherUseCase(getCurrentWeatherQuery *GetCurrentWeatherQuery, saveWeatherCommand *SaveWeatherCommand) *SaveCurrentWeatherUseCase {
	return &SaveCurrentWeatherUseCase{
		GetCurrentWeatherQuery: getCurrentWeatherQuery,
		SaveWeatherCommand:     saveWeatherCommand,
	}
}

func (u *SaveCurrentWeatherUseCase) Do() error {
	weather, err := u.GetCurrentWeatherQuery.Do()
	if err != nil {
		log.Println("Error getting current weather: ", err)
		return err
	}
	return u.SaveWeatherCommand.Do(weather)
}
