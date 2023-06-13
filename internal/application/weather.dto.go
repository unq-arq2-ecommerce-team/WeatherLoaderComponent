package application

import (
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type (
	WeatherDTO struct {
		City string  `json:"name"`
		Data DataDTO `json:"main"`
	}
	DataDTO struct {
		Temp     float64 `json:"temp"`
		FeelLike float64 `json:"feels_like"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	}
)

// IsInvalid : weatherDTO must have name (city) and some weather data
func (dto WeatherDTO) IsInvalid() bool {
	return dto.City == "" || dto.Data.IsInvalid()
}

func (dto DataDTO) IsInvalid() bool {
	return dto.Temp == 0 && dto.TempMin == 0 && dto.TempMax == 0 && dto.FeelLike == 0 && dto.Pressure == 0 && dto.Humidity == 0
}

func (dto WeatherDTO) String() string {
	return domain.ParseStruct(dto)
}
