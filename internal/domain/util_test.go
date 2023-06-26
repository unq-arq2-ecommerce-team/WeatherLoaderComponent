package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ParseStruct_WithCollectionsOrNested_Should_ReturnOk(t *testing.T) {
	anomStruct := struct {
		SomeArr    []string
		SomeMap    map[string]string
		SomeStruct struct{ SomeStr string }
	}{
		SomeArr:    []string{"a", "b"},
		SomeMap:    map[string]string{"a": "1", "b": "2", "c": "3"},
		SomeStruct: struct{ SomeStr string }{SomeStr: "sarasa"},
	}
	assert.Equal(t, `{"SomeArr":["a","b"],"SomeMap":{"a":"1","b":"2","c":"3"},"SomeStruct":{"SomeStr":"sarasa"}}`, ParseStruct(anomStruct))
}

func Test_ParseStruct_WithNoCommonEncodingJSON_Should_ReturnEmptyString(t *testing.T) {
	idString := func(x string) string { return x }
	assert.Equal(t, "", ParseStruct(idString))
}

func Test_GetDaysBetween(t *testing.T) {
	dateFrom1 := time.Date(2022, 10, 1, 12, 12, 12, 0, time.UTC)
	dateTo1 := time.Date(2023, 05, 5, 12, 12, 12, 0, time.UTC)
	dateFrom2 := time.Date(2023, 05, 1, 0, 0, 0, 0, time.UTC)
	dateTo2 := time.Date(2023, 05, 5, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, float64(216), GetDaysBetween(dateFrom1, dateTo1))
	assert.Equal(t, float64(216), GetDaysBetween(dateTo1, dateFrom1))
	assert.Equal(t, float64(4), GetDaysBetween(dateFrom2, dateTo2))
}
