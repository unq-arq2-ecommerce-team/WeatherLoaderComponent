package dto

import (
	"time"
)

type TemperatureAverageQuery struct {
	DateFrom time.Time `form:"dateFrom" binding:"required" time_format:"2006-01-02T15:04:05.000Z"`
	DateTo   time.Time `form:"dateTo" binding:"required" time_format:"2006-01-02T15:04:05.000Z"`
}

func (dto TemperatureAverageQuery) InvalidDates() bool {
	return !dto.DateTo.After(dto.DateFrom)
}

func (dto TemperatureAverageQuery) GetDateFrom() time.Time {
	return dto.DateFrom
}

func (dto TemperatureAverageQuery) GetDateTo() time.Time {
	return dto.DateTo
}
