package domain

import (
	"encoding/json"
)

func ParseStruct(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
