package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_AverageTemperature_Set(t *testing.T) {
	city, avgTemp := "Quilmes", 12.3332
	avgTempBase := &AverageTemperature{
		City:           city,
		AvgTemperature: avgTemp,
	}
	avgTempCopy := *avgTempBase
	dateFrom := time.Date(2023, 10, 1, 12, 12, 12, 0, time.UTC)
	dateTo := time.Date(2023, 10, 5, 12, 12, 12, 0, time.UTC)

	avgTempBase.Set(dateFrom, dateTo)

	assert.Equal(t, city, avgTempCopy.City)
	assert.Equal(t, avgTemp, avgTempCopy.AvgTemperature)
	assert.Zero(t, avgTempCopy.DateFrom)
	assert.Zero(t, avgTempCopy.DateTo)
	assert.Zero(t, avgTempCopy.DaysBetween)

	assert.Equal(t, city, avgTempBase.City)
	assert.Equal(t, avgTemp, avgTempBase.AvgTemperature)
	assert.Equal(t, dateFrom, avgTempBase.DateFrom)
	assert.Equal(t, dateTo, avgTempBase.DateTo)
	assert.Equal(t, float64(4), avgTempBase.DaysBetween)
}
