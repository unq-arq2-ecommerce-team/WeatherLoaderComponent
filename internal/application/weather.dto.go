package application

import (
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
)

type WeatherDTO struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		FeelLike float64 `json:"feels_like"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

func (dto WeatherDTO) String() string {
	return domain.ParseStruct(dto)
}
