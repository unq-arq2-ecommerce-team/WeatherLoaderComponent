package domain

import (
	"encoding/json"
	"math"
	"time"
)

func ParseStruct(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func GetDaysBetween(dateFrom, dateTo time.Time) float64 {
	return math.Abs(dateTo.Sub(dateFrom).Hours() / 24)
}
