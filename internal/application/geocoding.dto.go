package application

import "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"

type GeoCodingDTO struct {
	City string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func (dto GeoCodingDTO) String() string {
	return domain.ParseStruct(dto)
}
